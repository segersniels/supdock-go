package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	docker "supdock-go/src/docker"

	"github.com/segersniels/goutil"
)

func usage() {
	output := `Usage: supdock [options] [command]

Options:
	-h, --help         output usage information

Commands:
	stop              Stop a running container
	start             Start a stopped container
	logs              See the logs of a container
	rm                Remove a container
	rmi               Remove an image
	prune             Remove stopped containers and dangling images
	stats             See the stats of a container
	ssh               SSH into a container
	history           See the history of an image
	env               See the environment variables of a running container
	latest            Update to the latest version of supdock
`
	fmt.Print(output)
}

func help() {
	usage()
	dockerOut, err := exec.Command("docker", "--help").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", dockerOut)
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
		"prune",
		"history",
	}
	if util.Exists(commands, os.Args[1]) && len(os.Args) == 2 {
		docker.Execute(os.Args[1])
	} else {
		switch os.Args[1] {
		case "-h", "--help", "help":
			help()
		case "latest":
			util.Download("/usr/local/bin/supdock", "https://github.com/segersniels/supdock-go/raw/master/bin/supdock")
		default:
			docker.Standard(os.Args[1:])
		}
	}
}
