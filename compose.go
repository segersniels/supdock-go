package main

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	funk "github.com/thoas/go-funk"
)

type Compose struct {
	Name string
	Path string
}

func customCompose(args []string, path string) {
	cmd := exec.Command("docker-compose", append([]string{"-f", path}, args...)...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func checkIfLocalDockerComposeFile() (bool, string) {
	if _, err := os.Stat("./docker-compose.yaml"); !os.IsNotExist(err) {
		return true, "docker-compose.yaml"
	} else if _, err := os.Stat("./docker-compose.yml"); !os.IsNotExist(err) {
		return true, "docker-compose.yml"
	}
	return false, ""
}

func executeCompose(command string, question string) {
	var path string
	isLocal, file := checkIfLocalDockerComposeFile()
	if !isLocal {
		files := searchComposeFiles()
		names := funk.Map(files, func(c Compose) string {
			return c.Name
		})
		name := promptQuestion(question, names.([]string))
		project := funk.Find(files, func(c Compose) bool {
			return c.Name == name
		})
		path = project.(Compose).Path
	} else {
		path = file
	}
	switch command {
	case "up-detached":
		customCompose([]string{"up", "-d"}, path)
	default:
		customCompose([]string{command}, path)
	}
}

func passThroughCompose() {
	cmd := exec.Command("docker-compose", os.Args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func strip(path string) string {
	firstParent := filepath.Dir(path)
	secondParent := filepath.Base(filepath.Dir(firstParent))
	return secondParent + "/" + filepath.Base(firstParent)
}

func search(wg *sync.WaitGroup, root string, depth int, results chan Compose) {
	defer wg.Done()

	visit := func(path string, info os.FileInfo, err error) error {
		// When subdirectory is found and it isn't the root directory start a new parallel walk
		if info.IsDir() && path != root {
			if len(strings.Split(path, "/")) < depth {
				wg.Add(1)
				go search(wg, path, depth, results)
			}
			return filepath.SkipDir
		}
		if strings.Contains(path, "docker-compose.yaml") || strings.Contains(path, "docker-compose.yml") {
			results <- Compose{
				Name: strip(path),
				Path: path,
			}
		}
		return nil
	}

	err := filepath.Walk(root, visit)
	if err != nil {
		log.Fatal(err)
	}
}

func searchComposeFiles() []Compose {
	var wg sync.WaitGroup
	var files []Compose
	results := make(chan Compose)

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	search(&wg, usr.HomeDir, depth, results)
	go func() {
		for file := range results {
			files = append(files, file)
		}
	}()
	wg.Wait()

	return files
}
