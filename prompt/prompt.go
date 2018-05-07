package prompt

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/segersniels/goutil"
)

// Docker : standard docker execution
func Docker(args []string) {
	var errbuf bytes.Buffer
	cmd := exec.Command("docker", args...)
	cmd.Stderr = &errbuf
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		util.Error(strings.TrimSpace(errbuf.String()))
	}
}

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

// Exec : specify the execution
func Exec(command string, ids []string, names []string, question string) {
	if len(ids) > 1 && len(names) > 1 {
		options := constructChoices(ids, names)
		answer := util.Question(question, options)
		switch command {
		case "ssh":
			shell := util.Question("Which shell is the container using?", []string{"bash", "ash"})
			Docker([]string{"exec", "-ti", strings.Split(answer, " - ")[0], shell})
		case "env":
			Docker([]string{"exec", "-ti", strings.Split(answer, " - ")[0], "env"})
		default:
			Docker([]string{command, strings.Split(answer, " - ")[0]})
		}
	} else {
		util.Warn("No options found to construct prompt")
	}
}
