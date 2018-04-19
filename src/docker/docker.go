package docker

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/segersniels/goutil"
)

func fullCommandExecute(command string) string {
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stderr = &errbuf
	cmd.Stdout = &outbuf
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		util.Error(strings.TrimSpace(errbuf.String()))
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
		util.Error(strings.TrimSpace(errbuf.String()))
	}
}

// Execute : return a prompt and execute
func Execute(command string) {
	psIds := strings.Split(fullCommandExecute("docker ps |tail -n +2 |awk '{print $1}'"), "\n")
	psaIds := strings.Split(fullCommandExecute("docker ps -a |tail -n +2 |awk '{print $1}'"), "\n")
	imageIds := strings.Split(fullCommandExecute("docker images |tail -n +2 |awk '{print $3}'"), "\n")
	psNames := strings.Split(fullCommandExecute("docker ps |tail -n +2 |awk '{print $NF}'"), "\n")
	psaNames := strings.Split(fullCommandExecute("docker ps -a |tail -n +2 |awk '{print $NF}'"), "\n")
	imageNames := strings.Split(fullCommandExecute("docker images |tail -n +2 |awk '{print $1}'"), "\n")
	switch command {
	case "logs":
		constructPrompt("logs", psaIds, psaNames, "Which container would you like to see the logs of?")
	case "start":
		constructPrompt("start", psaIds, psaNames, "Which container would you like to start?")
	case "restart":
		constructPrompt("restart", psIds, psNames, "Which container would you like to restart?")
	case "stop":
		constructPrompt("stop", psIds, psNames, "Which container would you like to stop?")
	case "ssh":
		constructPrompt("ssh", psIds, psNames, "Which container would you like to connect with?")
	case "env":
		constructPrompt("env", psIds, psNames, "Which container would you like to see the environment variables of?")
	case "rm":
		constructPrompt("rm", psaIds, psaNames, "Which container would you like to remove?")
	case "rmi":
		constructPrompt("rmi", imageIds, imageNames, "Which image would you like to remove?")
	case "history":
		constructPrompt("history", imageIds, imageNames, "Which image would you like to see the history of?")
	case "stats":
		constructPrompt("stats", psIds, psNames, "Which container would you like to see that stats of?")
	}
}
