package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"supdock-go/src/docker"

	"github.com/segersniels/goutil"
	"github.com/urfave/cli"
)

func update() {
	version := strings.TrimSpace(util.ExecuteWithOutput("curl --silent 'https://api.github.com/repos/segersniels/supdock-go/releases/latest' |grep tag_name |awk '{print $2}' |tr -d '\",v'"))
	distro := strings.TrimSpace(runtime.GOOS)
	if distro != "darwin" && distro != "linux" {
		fmt.Println("ERR: operating system does not equal either linux or darwin")
		os.Exit(0)
	}
	fmt.Println("Updating to version", version+"-"+distro)
	util.Download("/usr/local/bin/supdock", "https://github.com/segersniels/supdock-go/releases/download/v"+version+"/supdock"+version+"_"+distro+"_amd64")
}

func main() {
	app := cli.NewApp()
	app.Name = "supdock"
	app.Usage = "What's Up, Dock(er)?"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		{
			Name:  "stop",
			Usage: "Stop a running container",
			Action: func(c *cli.Context) error {
				docker.Execute(os.Args[1:])
				return nil
			},
		},
		{
			Name:  "start",
			Usage: "Start a stopped container",
			Action: func(c *cli.Context) error {
				docker.Execute(os.Args[1:])
				return nil
			},
		},
		{
			Name:  "logs",
			Usage: "See the logs of a container",
			Action: func(c *cli.Context) error {
				docker.Execute(os.Args[1:])
				return nil
			},
		},
		{
			Name:  "rm",
			Usage: "Remove a container",
			Action: func(c *cli.Context) error {
				docker.Execute(os.Args[1:])
				return nil
			},
		},
		{
			Name:  "rmi",
			Usage: "Remove an image",
			Action: func(c *cli.Context) error {
				docker.Execute(os.Args[1:])
				return nil
			},
		},
		{
			Name:  "prune",
			Usage: "Remove stopped containers and dangling images",
			Action: func(c *cli.Context) error {
				docker.Standard([]string{"system", os.Args[1], "-f"})
				return nil
			},
		},
		{
			Name:  "stats",
			Usage: "See the stats of a container",
			Action: func(c *cli.Context) error {
				docker.Execute(os.Args[1:])
				return nil
			},
		},
		{
			Name:  "ssh",
			Usage: "SSH into a container",
			Action: func(c *cli.Context) error {
				docker.Execute(os.Args[1:])
				return nil
			},
		},
		{
			Name:  "history",
			Usage: "See the history of an image",
			Action: func(c *cli.Context) error {
				docker.Execute(os.Args[1:])
				return nil
			},
		},
		{
			Name:  "env",
			Usage: "See the environment variables of a running container",
			Action: func(c *cli.Context) error {
				docker.Execute(os.Args[1:])
				return nil
			},
		},
		{
			Name:    "update",
			Aliases: []string{"latest"},
			Flags: []cli.Flag{
				cli.BoolFlag{},
			},
			Usage: "Update to the latest version",
			Action: func(c *cli.Context) error {
				update()
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
