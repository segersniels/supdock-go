package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/segersniels/goutil"
	"github.com/segersniels/supdock-go/prompt"
)

var psIds, psaIds, imageIds, psNames, psaNames, imageNames []string

func initialise() {
	ids, _ := util.ExecuteWithOutput("docker ps -q")
	psIds = strings.Split(ids, "\n")
	ids, _ = util.ExecuteWithOutput("docker ps -aq")
	psaIds = strings.Split(ids, "\n")
	ids, _ = util.ExecuteWithOutput("docker images -q")
	imageIds = strings.Split(ids, "\n")

	names, _ := util.ExecuteWithOutput("docker ps |tail -n +2 |awk '{print $NF}'")
	psNames = strings.Split(names, "\n")
	names, _ = util.ExecuteWithOutput("docker ps -a |tail -n +2 |awk '{print $NF}'")
	psaNames = strings.Split(names, "\n")
	names, _ = util.ExecuteWithOutput("docker images |tail -n +2 |awk '{print $1}'")
	imageNames = strings.Split(names, "\n")
}

func usage() {
	output := `Usage: supdock [options] [command]

Options:
	-h, --help         output usage information

Commands:
	stop              Stop a running container
	destroy           Stop all running containers
	start             Start a stopped container
	restart           Restart a running container
	logs              See the logs of a container
	rm                Remove a container
	rmi               Remove an image
	prune             Remove stopped containers and dangling images
	stats             See the stats of a container
	ssh               SSH into a container
	history           See the history of an image
	history           Inspect a container
	env               See the environment variables of a running container
	latest, update    Update to the latest version of supdock
`
	fmt.Println(output)
}

func version() {
	app := "supdock"
	version := "0.1.2-rc.2"
	fmt.Println(app, "version", version)
}

func help() {
	usage()
	dockerOut, err := exec.Command("docker", "--help").Output()
	if err != nil {
		util.Error(err)
	}
	fmt.Printf("%s", dockerOut)
}

func update() {
	output, _ := util.ExecuteWithOutput("curl --silent 'https://api.github.com/repos/segersniels/supdock-go/releases/latest' |grep tag_name |awk '{print $2}' |tr -d '\",v'")
	version := strings.TrimSpace(output)
	distro := strings.TrimSpace(runtime.GOOS)
	if distro != "darwin" && distro != "linux" {
		util.Error("Operating system does not equal linux or darwin")
	}
	fmt.Println("Updating to version", version+"-"+distro)
	err := util.Download("/usr/local/bin/supdock", "https://github.com/segersniels/supdock-go/releases/download/v"+version+"/supdock_"+version+"_"+distro+"_amd64")
	if err != nil {
		util.Error(err)
	}
}

func execute(command string) {
	initialise()
	switch command {
	case "logs":
		prompt.Exec("logs", psaIds, psaNames, "Which container would you like to see the logs of?")
	case "start":
		prompt.Exec("start", psaIds, psaNames, "Which container would you like to start?")
	case "restart":
		prompt.Exec("restart", psIds, psNames, "Which container would you like to restart?")
	case "stop":
		prompt.Exec("stop", psIds, psNames, "Which container would you like to stop?")
	case "ssh":
		prompt.Exec("ssh", psIds, psNames, "Which container would you like to connect with?")
	case "env":
		prompt.Exec("env", psIds, psNames, "Which container would you like to see the environment variables of?")
	case "rm":
		prompt.Exec("rm", psaIds, psaNames, "Which container would you like to remove?")
	case "rmi":
		prompt.Exec("rmi", imageIds, imageNames, "Which image would you like to remove?")
	case "history":
		prompt.Exec("history", imageIds, imageNames, "Which image would you like to see the history of?")
	case "stats":
		prompt.Exec("stats", psIds, psNames, "Which container would you like to see that stats of?")
	case "inspect":
		prompt.Exec("inspect", psIds, psNames, "Which container would you like to inspect?")
	}
}

func main() {
	if len(os.Args) < 2 {
		help()
		os.Exit(0)
	}
	commands := []string{
		"logs",
		"start",
		"stop",
		"rm",
		"rmi",
		"ssh",
		"stats",
		"env",
		"history",
		"restart",
		"inspect",
	}
	if util.Exists(commands, os.Args[1]) && len(os.Args) == 2 {
		execute(os.Args[1])
	} else {
		switch os.Args[1] {
		case "-h", "--help", "help":
			help()
		case "-v", "--version", "version":
			version()
		case "latest", "update":
			update()
		case "prune":
			err := util.Execute("docker system prune -f", []string{})
			if err != nil {
				util.Error(err)
			}
		case "destroy":
			err := util.Execute("docker stop $(docker ps -q)", []string{})
			if err != nil {
				util.Error(err)
			}
		default:
			cmd := exec.Command("docker", os.Args[1:]...)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			err := cmd.Run()
			if err != nil {
				util.Error(err)
			}
		}
	}
}
