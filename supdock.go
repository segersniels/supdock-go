package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	docker "./lib"
)

func help() {
	supdockOut, err := exec.Command("supdock", "-h").Output()
	if err != nil {
		log.Fatal(err)
	}
	dockerOut, err := exec.Command("docker", "-h").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", supdockOut)
	fmt.Printf("%s", dockerOut)
}

func main() {
	if len(os.Args) < 2 {
		help()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "logs":
		docker.Execute("logs")
	case "start":
		docker.Execute("start")
	case "stop":
		docker.Execute("stop")
	case "ssh":
		docker.Execute("ssh")
	default:
		docker.Standard(os.Args[1:])
	}
}
