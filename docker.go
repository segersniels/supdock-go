package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
)

func start(id string) {
	err := docker.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}

func stop(id string) {
	err := docker.ContainerStop(context.Background(), id, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}

func stopParallel(id string, wg *sync.WaitGroup) {
	defer wg.Done()
	err := docker.ContainerStop(context.Background(), id, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}

func restartParallel(id string, wg *sync.WaitGroup) {
	defer wg.Done()
	err := docker.ContainerRestart(context.Background(), id, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}

func restart(id string) {
	err := docker.ContainerRestart(context.Background(), id, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}

func removeParallel(removeType string, id string, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	switch removeType {
	case "container":
		err = docker.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{})
	case "image":
		_, err = docker.ImageRemove(context.Background(), id, types.ImageRemoveOptions{})
	case "image-force":
		_, err = docker.ImageRemove(context.Background(), id, types.ImageRemoveOptions{Force: true})
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}

func remove(removeType string, id string) {
	var err error
	switch removeType {
	case "container":
		err = docker.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{})
	case "image":
		_, err = docker.ImageRemove(context.Background(), id, types.ImageRemoveOptions{})
	case "image-force":
		_, err = docker.ImageRemove(context.Background(), id, types.ImageRemoveOptions{Force: true})
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}

func passThroughDocker() {
	cmd := exec.Command("docker", os.Args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func customDocker(args []string) {
	var errbuf bytes.Buffer
	if exists(args, "|") {
		index := getIndex(args, "|")
		length := len(args)
		cmd := exec.Command("docker", args[0:index]...)
		pipeCmd := exec.Command(args[index+1], args[index+2:length]...)

		pipeCmd.Stdin, _ = cmd.StdoutPipe()
		pipeCmd.Stdout = os.Stdout
		pipeCmd.Stderr = &errbuf

		err := pipeCmd.Start()
		if err != nil {
			log.Fatal(strings.TrimSpace(errbuf.String()))
		}

		err = cmd.Run()
		if err != nil {
			log.Fatal(strings.TrimSpace(errbuf.String()))
		}

		err = pipeCmd.Wait()
		if err != nil {
			log.Fatal(strings.TrimSpace(errbuf.String()))
		}
	} else {
		cmd := exec.Command("docker", args...)
		cmd.Stderr = &errbuf
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		err := cmd.Run()
		if err != nil {
			log.Fatal(strings.TrimSpace(errbuf.String()))
		}
	}
}

func executeDocker(command string, ids []string, names []string, question string) {
	if len(ids) >= 1 && len(names) >= 1 {
		id := selectID(ids, names, question)
		switch command {
		case "ssh":
			shell := promptQuestion("Which shell is the container using?", []string{"bash", "ash"})
			customDocker([]string{"exec", "-ti", id, shell})
		case "env":
			customDocker([]string{"exec", "-ti", id, "env"})
		case "logs-force":
			customDocker([]string{"logs", "-f", id})
		case "stats-no-stream":
			customDocker([]string{"stats", "--no-stream", id})
		default:
			customDocker([]string{command, id})
		}
	} else {
		log.Fatal("No options found to construct prompt")
	}
}
