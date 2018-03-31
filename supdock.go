package main

import (
	"os"

	docker "./src/docker"
	util "./src/util"
)

func main() {
	if len(os.Args) < 2 {
		util.Help()
		os.Exit(0)
	}
	commands := []string{"logs", "start", "stop", "rm", "rmi", "ssh", "stats", "env", "prune", "history"}
	if util.SliceExists(commands, os.Args[1]) && len(os.Args) == 2 {
		docker.Execute(os.Args[1])
	} else {
		docker.Standard(os.Args[1:])
	}
}
