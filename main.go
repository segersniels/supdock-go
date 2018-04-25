package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"supdock-go/src/docker"

	"github.com/segersniels/goutil"
)

func usage() {
	output := `Usage: supdock [options] [command]

Options:
	-h, --help         output usage information

Commands:
	stop              Stop a running container
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
	version := "0.1.2"
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
	version := strings.TrimSpace(util.ExecuteWithOutput("curl --silent 'https://api.github.com/repos/segersniels/supdock-go/releases/latest' |grep tag_name |awk '{print $2}' |tr -d '\",v'"))
	distro := strings.TrimSpace(runtime.GOOS)
	if distro != "darwin" && distro != "linux" {
		util.Error("Operating system does not equal linux or darwin")
	}
	fmt.Println("Updating to version", version+"-"+distro)
	util.Download("/usr/local/bin/supdock", "https://github.com/segersniels/supdock-go/releases/download/v"+version+"/supdock_"+version+"_"+distro+"_amd64")
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
		docker.Execute(os.Args[1])
	} else {
		switch os.Args[1] {
		case "-h", "--help", "help":
			help()
		case "-v", "--version", "version":
			version()
		case "latest", "update":
			update()
		case "prune":
			docker.Standard([]string{"system", "prune", "-f"})
		default:
			docker.Standard(os.Args[1:])
		}
	}
}
