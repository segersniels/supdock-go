package command

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/segersniels/goutil"
)

// PassThrough : passthrough docker execution
func PassThrough() {
	cmd := exec.Command("docker", os.Args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		util.Error(err)
	}
}

func docker(args []string) {
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

func execute(command string, ids []string, names []string, question string) {
	if len(ids) > 1 && len(names) > 1 {
		options := constructChoices(ids, names)
		answer := util.Question(question, options)
		switch command {
		case "ssh":
			shell := util.Question("Which shell is the container using?", []string{"bash", "ash"})
			docker([]string{"exec", "-ti", strings.Split(answer, " - ")[0], shell})
		case "env":
			docker([]string{"exec", "-ti", strings.Split(answer, " - ")[0], "env"})
		case "logs -f":
			docker([]string{"logs", "-f", strings.Split(answer, " - ")[0]})
		default:
			docker([]string{command, strings.Split(answer, " - ")[0]})
		}
	} else {
		util.Warn("No options found to construct prompt")
	}
}
