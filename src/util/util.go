package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/AlecAivazis/survey"
)

// SliceExists : check if a slice contains a specific value
func SliceExists(slice interface{}, item interface{}) bool {
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

// AskQuestion : ask the user a question prompt using survey package
func AskQuestion(question string, options []string) string {
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

// Repo : repo object
type Repo struct {
	Repo  string `json:"repo"`
	Short string `json:"short"`
}

// ParseJson : read through a json file
func ParseJson(filename string) []Repo {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	var repos []Repo
	json.Unmarshal(byteValue, &repos)
	return repos
}
