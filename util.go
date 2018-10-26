package main

import (
	"reflect"
	"strings"

	"github.com/AlecAivazis/survey"
	log "github.com/sirupsen/logrus"
)

func constructChoices(ids []string, names []string) []string {
	choices := []string{}
	for index, id := range ids {
		choice := id + " - " + strings.TrimLeft(names[index], "/")
		if choice != " - " {
			choices = append(choices, choice)
		}
	}
	return choices
}

func selectID(ids []string, names []string, question string) string {
	options := constructChoices(ids, names)
	answer := promptQuestion(question, options)
	id := strings.Split(answer, " - ")[0]
	return id
}

func getIndex(slice interface{}, item interface{}) int {
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

func exists(slice interface{}, item interface{}) bool {
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

func promptQuestion(question string, options []string) string {
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
		log.Fatal(err)
	}
	return answers.Selection
}
