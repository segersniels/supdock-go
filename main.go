package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/segersniels/goutil"
	"github.com/segersniels/supdock-go/prompt"
	"github.com/urfave/cli"
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

func docker() {
	cmd := exec.Command("docker", os.Args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		util.Error(err)
	}
}

func main() {
	utilities := []string{"-h", "--help", "-v", "--version"}
	commands := []string{
		"latest",
		"upgrade",
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
		"prune",
		"destroy",
		"shutdown",
		"memory",
	}

	if len(os.Args) > 1 {
		if util.Exists(commands, os.Args[1]) || util.Exists(utilities, os.Args[1]) {
			if !util.Exists(utilities, os.Args[1]) {
				initialise()
			}
			app := cli.NewApp()
			app.Name = "supdock"
			app.Usage = "What's Up Dock(er)?"
			app.Version = "0.1.5"
			app.Commands = []cli.Command{
				{
					Name:  "logs",
					Usage: "See the logs of a container",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "details",
							Usage: "Show extra details provided to logs",
						},
						cli.BoolFlag{
							Name:  "f, follow",
							Usage: "Follow log output",
						},
						cli.StringFlag{
							Name:  "since",
							Usage: "Show logs since timestamp",
						},
						cli.StringFlag{
							Name:  "tail",
							Usage: "Number of lines to show from the end of the logs",
						},
						cli.BoolFlag{
							Name:  "t, timestamps",
							Usage: "Show timestamps",
						},
						cli.StringFlag{
							Name:  "until",
							Usage: "Show logs before a timestamp",
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 {
							if c.NumFlags() == 2 && c.Bool("f") {
								prompt.Exec("logs -f", psaIds, psaNames, "Which container would you like to see the logs of?")
							} else {
								prompt.Exec("logs", psaIds, psaNames, "Which container would you like to see the logs of?")
							}
						} else {
							docker()
						}
						return nil
					},
				},
				{
					Name:  "start",
					Usage: "Start a stopped container",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "a, attach",
							Usage: "Attach STDOUT/STDERR and forward signals",
						},
						cli.StringFlag{
							Name:  "checkpoint",
							Usage: "Restore from this checkpoint",
						},
						cli.StringFlag{
							Name:  "checkpoint-dir",
							Usage: "Use a custom checkpoint storage directory",
						},
						cli.StringFlag{
							Name:  "detach-keys",
							Usage: "Override the key sequence for detaching",
						},
						cli.BoolFlag{
							Name:  "i, interactive",
							Usage: "Attach container's STDIN",
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 && c.NumFlags() == 0 {
							prompt.Exec("start", psaIds, psaNames, "Which container would you like to start?")
						} else {
							docker()
						}
						return nil
					},
				},
				{
					Name:  "restart",
					Usage: "Restart a running container",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "t, time",
							Usage: "Seconds to wait for stop before killing the container",
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 && c.NumFlags() == 0 {
							prompt.Exec("restart", psIds, psNames, "Which container would you like to restart?")
						} else {
							docker()
						}
						return nil
					},
				},
				{
					Name:  "stop",
					Usage: "Stop a running container",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "t, time",
							Usage: "Seconds to wait for stop before killing the container",
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 && c.NumFlags() == 0 {
							prompt.Exec("stop", psIds, psNames, "Which container would you like to stop?")
						} else {
							docker()
						}
						return nil
					},
				},
				{
					Name:  "ssh",
					Usage: "SSH into a container",
					Action: func(c *cli.Context) error {
						prompt.Exec("ssh", psIds, psNames, "Which container would you like to connect with?")
						return nil
					},
				},
				{
					Name:  "env",
					Usage: "See the environment variables of a running container",
					Action: func(c *cli.Context) error {
						prompt.Exec("env", psIds, psNames, "Which container would you like to see the environment variables of?")
						return nil
					},
				},
				{
					Name:  "rm",
					Usage: "Remove a container",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "f, force",
							Usage: "Force the removal of a running container (uses SIGKILL)",
						},
						cli.BoolFlag{
							Name:  "l, link",
							Usage: "Remove the specified link",
						},
						cli.BoolFlag{
							Name:  "v, volumes",
							Usage: "Remove the volumes associated with the container",
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 && c.NumFlags() == 0 {
							prompt.Exec("rm", psaIds, psaNames, "Which container would you like to remove?")
						} else {
							docker()
						}
						return nil
					},
				},
				{
					Name:  "rmi",
					Usage: "Remove an image",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "f, force",
							Usage: "Force removal of the image",
						},
						cli.BoolFlag{
							Name:  "no-prune",
							Usage: "Do not delete untagged parents",
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 && c.NumFlags() == 0 {
							prompt.Exec("rmi", imageIds, imageNames, "Which image would you like to remove?")
						} else {
							docker()
						}
						return nil
					},
				},
				{
					Name:  "history",
					Usage: "See the history of an image",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "format",
							Usage: "Pretty-print images using a Go template",
						},
						cli.BoolFlag{
							Name:  "H, human",
							Usage: "Print sizes and dates in human readable format",
						},
						cli.BoolFlag{
							Name:  "no-trunc",
							Usage: "Don't truncate output",
						},
						cli.BoolFlag{
							Name:  "q, quiet",
							Usage: "Only show numeric IDs",
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 && c.NumFlags() == 0 {
							prompt.Exec("history", imageIds, imageNames, "Which image would you like to see the history of?")
						} else {
							docker()
						}
						return nil
					},
				},
				{
					Name:  "stats",
					Usage: "See the stats of a container",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "a, all",
							Usage: "Show all containers (default shows just running)",
						},
						cli.StringFlag{
							Name:  "format",
							Usage: "Pretty-print images using a Go template",
						},
						cli.BoolFlag{
							Name:  "no-stream",
							Usage: "Disable streaming stats and only pull the first",
						},
						cli.BoolFlag{
							Name:  "no-trunc",
							Usage: "Do not truncate output",
						},
						cli.BoolFlag{
							Name:  "s, select",
							Usage: "Select the container from a select prompt",
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 && c.Bool("s") {
							prompt.Exec("stats", psIds, psNames, "Which container would you like to see that stats of?")
						} else {
							docker()
						}
						return nil
					},
				},
				{
					Name:  "inspect",
					Usage: "Inspect a container",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "f, format",
							Usage: "Format the output using the given Go template",
						},
						cli.BoolFlag{
							Name:  "s, size",
							Usage: "Display total file sizes if the type is container",
						},
						cli.StringFlag{
							Name:  "type",
							Usage: "Return JSON for specified type",
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 && c.NumFlags() == 0 {
							prompt.Exec("inspect", psIds, psNames, "Which container would you like to inspect?")
						} else {
							docker()
						}
						return nil
					},
				},
				{
					Name:  "prune",
					Usage: "Remove stopped containers and dangling images",
					Action: func(c *cli.Context) error {
						err := util.Execute("docker system prune -f", []string{})
						if err != nil {
							util.Error(err)
						}
						return nil
					},
				},
				{
					Name:    "destroy",
					Usage:   "Stop all running containers",
					Aliases: []string{"shutdown"},
					Action: func(c *cli.Context) error {
						err := util.Execute("docker stop $(docker ps -q)", []string{})
						if err != nil {
							util.Error(err)
						}
						return nil
					},
				},
				{
					Name:  "memory",
					Usage: "See the memory usage of all running containers",
					Action: func(c *cli.Context) error {
						err := util.Execute("docker ps -q | xargs  docker stats --no-stream", []string{})
						if err != nil {
							util.Error(err)
						}
						return nil
					},
				},
				{
					Name:    "latest",
					Usage:   "Update to the latest version of supdock",
					Aliases: []string{"upgrade"},
					Action: func(c *cli.Context) error {
						update()
						return nil
					},
				},
			}
			err := app.Run(os.Args)
			if err != nil {
				util.Warn(err)
				os.Exit(0)
			}
		} else {
			docker()
		}
	} else {
		docker()
	}
}
