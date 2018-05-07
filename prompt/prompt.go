package prompt

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/segersniels/goutil"
)

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

// Exec : specify the execution
func Exec(command string, ids []string, names []string, question string) {
	if len(ids) > 1 && len(names) > 1 {
		options := constructChoices(ids, names)
		answer, err := util.Question(question, options)
		if err != nil {
			util.Error(err)
		}
		switch command {
		case "ssh":
			shell, err := util.Question("Which shell is the container using?", []string{"bash", "ash"})
			if err != nil {
				util.Error(err)
			}
			docker([]string{"exec", "-ti", strings.Split(answer, " - ")[0], shell})
		case "env":
			docker([]string{"exec", "-ti", strings.Split(answer, " - ")[0], "env"})
		default:
			docker([]string{command, strings.Split(answer, " - ")[0]})
		}
	} else {
		util.Warn("No options found to construct prompt")
	}
}
