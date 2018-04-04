package util

import (
	"fmt"
	"os"
	"reflect"

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
