package util

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/AlecAivazis/survey"
)

// Input : ask for user input
func Input(question string) (string, error) {
	var qs = []*survey.Question{
		{
			Name: "selection",
			Prompt: &survey.Input{
				Message: question,
			},
		},
	}
	answers := struct {
		Selection string
	}{}
	err := survey.Ask(qs, &answers)
	return answers.Selection, err
}

// InputDefault : ask for user input with default value
func InputDefault(question string, value string) (string, error) {
	var qs = []*survey.Question{
		{
			Name: "selection",
			Prompt: &survey.Input{
				Message: question,
				Default: value,
			},
		},
	}
	answers := struct {
		Selection string
	}{}
	err := survey.Ask(qs, &answers)
	return answers.Selection, err
}

// Question : ask the user a question prompt using survey package
func Question(question string, options []string) (string, error) {
	var qs = []*survey.Question{
		{
			Name: "selection",
			Prompt: &survey.Select{
				Message: question,
				Options: options,
			},
		},
	}
	answers := struct {
		Selection string
	}{}
	err := survey.Ask(qs, &answers)
	return answers.Selection, err
}

// Search : search object needed for SearchForFile
type Search struct {
	path string
	name string
}

func scan(wg *sync.WaitGroup, folder string, depth int, results *[]Search, name string) {
	defer wg.Done()

	files, _ := ioutil.ReadDir(folder)
	var directories []string

	for _, file := range files {
		path := folder + "/" + file.Name()

		if file.IsDir() {
			directories = append(directories, path)
			continue
		}

		if !file.IsDir() && file.Name() == name {
			*results = append(*results, Search{
				path: filepath.Dir(path),
				name: filepath.Base(filepath.Dir(path)),
			})
			return
		}
	}

	if depth > 1 {
		for _, folder := range directories {
			wg.Add(1)
			go scan(wg, folder, depth-1, results, name)
		}
	}

	return
}

// SearchForFile : search home directory for a file and return all the results in an object containing path and name
func SearchForFile(name string) []Search {
	usr, _ := user.Current()
	searches := []Search{}
	var wg = new(sync.WaitGroup)
	wg.Add(1)
	go scan(wg, usr.HomeDir, 5, &searches, name)
	wg.Wait()
	return searches
}

// Download : download a given file to a location
func Download(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// Exists : check if a slice contains a specific value
func Exists(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("SliceExists() given a non-slice type")
	}
	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}
	return false
}

// GetIndex : search a slice for a value and return its index
func GetIndex(slice interface{}, item interface{}) int {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("SliceExists() given a non-slice type")
	}
	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return i
		}
	}
	return -1
}

// Min : return the minimum value of a slice
func Min(values []int) int {
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

// Max : return the maximum value of a slice
func Max(values []int) int {
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

// ExecuteWithOutput : execute a command through bash -c and return the stdout
func ExecuteWithOutput(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	var outbuf, errbuf bytes.Buffer
	cmd.Stderr = &errbuf
	cmd.Stdout = &outbuf
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return outbuf.String(), err
}

// Execute : execute a command through bash -c, pass []string{} as env to ignore custom vars
func Execute(command string, env []string) error {
	var errbuf bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stderr = &errbuf
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if len(env) > 0 {
		cmd.Env = os.Environ()
		for _, e := range env {
			cmd.Env = append(cmd.Env, e)
		}
	}
	err := cmd.Run()
	return err
}
