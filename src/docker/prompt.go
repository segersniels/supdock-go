package docker

import (
	"strings"

	"github.com/segersniels/goutil"
)

func constructChoices(ids []string, names []string) []string {
	choices := []string{}
	for index, id := range ids {
		choice := id + " - " + names[index]
		if choice != " - " {
			choices = append(choices, choice)
		}
	}
	return choices
}

func constructPrompt(command string, ids []string, names []string, question string) {
	if len(ids) > 1 && len(names) > 1 {
		options := constructChoices(ids, names)
		answer := util.Question(question, options)
		switch command {
		case "ssh":
			shell := util.Question("Which shell is the container using?", []string{"bash", "ash"})
			Standard([]string{"exec", "-ti", strings.Split(answer, " - ")[0], shell})
		case "env":
			Standard([]string{"exec", "-ti", strings.Split(answer, " - ")[0], "env"})
		default:
			Standard([]string{command, strings.Split(answer, " - ")[0]})
		}
	} else {
		util.Warn("No options found to construct prompt")
	}
}
