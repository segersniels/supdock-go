package util

import (
	"bytes"
	"fmt"
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

// RequireInput : ask for input to the user
func RequireInput(question string) string {
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
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	return answers.Selection
}

// Question : ask the user a question prompt using survey package
func Question(question string, options []string) string {
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
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	return answers.Selection
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
func Download(filepath string, url string) (err error) {
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// GetIndexString : search a slice for a value and return its index
func GetIndexString(value string, slice []string) int {
	for index, v := range slice {
		if value == v {
			return index
		}
	}
	return -1
}

// GetIndexInt : search a slice for a value and return its index
func GetIndexInt(value int, slice []int) int {
	for index, v := range slice {
		if value == v {
			return index
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
func ExecuteWithOutput(command string) string {
	cmd := exec.Command("bash", "-c", command)
	var outbuf, errbuf bytes.Buffer
	cmd.Stderr = &errbuf
	cmd.Stdout = &outbuf
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Print(errbuf.String())
	}
	return outbuf.String()
}

// Execute : execute a command through bash -c, pass []string{} as env to ignore custom vars
func Execute(command string, env []string) {
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
	if err != nil {
		fmt.Print(errbuf.String())
		os.Exit(0)
	}
}
