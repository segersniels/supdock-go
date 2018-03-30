package docker

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey"
)

func fullCommandExecute(command string) string {
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stderr = &errbuf
	cmd.Stdout = &outbuf
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Print(errbuf.String())
	}
	return outbuf.String()
}

// Standard : default execution of docker
func Standard(args []string) {
	var errbuf bytes.Buffer
	cmd := exec.Command("docker", args...)
	cmd.Stderr = &errbuf
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Print(errbuf.String())
	}
}

func constructChoices(ids []string, names []string) []string {
	var choices = []string{}
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
			return
		}
		var cmd = []string{
			command,
			strings.Split(answers.Selection, " - ")[0],
		}
		if command != "ssh" {
			Standard(cmd)
		} else {
			command := strings.Split("exec -ti "+strings.Split(answers.Selection, " - ")[0]+" bash", " ")
			Standard(command)
		}
	} else {
		fmt.Print("ERR: No options found to construct prompt")
	}
}

// Execute : return a prompt and execute
func Execute(command string) {
	psIds := strings.Split(fullCommandExecute("docker ps |tail -n +2 |awk '{print $1}'"), "\n")
	psNames := strings.Split(fullCommandExecute("docker ps |tail -n +2 |awk '{print $NF}'"), "\n")
	psaIds := strings.Split(fullCommandExecute("docker ps -a |tail -n +2 |awk '{print $1}'"), "\n")
	psaNames := strings.Split(fullCommandExecute("docker ps -a |tail -n +2 |awk '{print $NF}'"), "\n")
	// imageIds := fullCommandExecute("docker images |tail -n +2 |awk '{print $3}'")
	// imageNames := fullCommandExecute("docker images |tail -n +2 |awk '{print $1}'")
	switch command {
	case "logs":
		constructPrompt("logs", psaIds, psaNames, "Which container would you like to see the logs of?")
	case "start":
		constructPrompt("start", psaIds, psaNames, "Which container would you like to start?")
	case "stop":
		constructPrompt("stop", psIds, psNames, "Which container would you like to stop?")
	case "ssh":
		constructPrompt("ssh", psIds, psNames, "Which container would you like to connect with?")
	}
}
