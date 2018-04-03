package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	docker "./src/docker"
	util "./src/util"
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
`
	fmt.Print(output)
}

// Help : call docker help as output for supdock
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
	commands := []string{"logs", "start", "stop", "rm", "rmi", "ssh", "stats", "env", "prune", "history"}
	if util.SliceExists(commands, os.Args[1]) && len(os.Args) == 2 {
		docker.Execute(os.Args[1])
	} else {
		if os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "help" {
			help()
		} else {
			docker.Standard(os.Args[1:])
		}
	}
}
