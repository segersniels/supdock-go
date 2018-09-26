package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	util "github.com/segersniels/goutil"
	"github.com/urfave/cli"
)

func passThroughDocker() {
	cmd := exec.Command("docker", os.Args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		util.Error(err)
	}
}

func customDocker(args []string) {
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

func extractNames(commands []cli.Command) []string {
	var names []string
	for _, command := range commands {
		names = append(names, command.Name)
		for _, alias := range command.Aliases {
			names = append(names, alias)
		}
	}
	return names
}

func constructChoices(ids []string, names []string) []string {
	choices := []string{}
	for index, id := range ids {
		choice := id + " - " + names[index]
		if choice != " - " {
			choices = append(choices, choice)
		}
	}
	return choices
}

func execute(command string, ids []string, names []string, question string) {
	if len(ids) >= 1 && len(names) >= 1 {
		options := constructChoices(ids, names)
		answer := util.Question(question, options)
		id := strings.Split(answer, " - ")[0]
		switch command {
		case "ssh":
			shell := util.Question("Which shell is the container using?", []string{"bash", "ash"})
			customDocker([]string{"exec", "-ti", id, shell})
		case "env":
			customDocker([]string{"exec", "-ti", id, "env"})
		case "logs -f":
			customDocker([]string{"logs", "-f", id})
		default:
			customDocker([]string{command, id})
		}
	} else {
		util.Warn("No options found to construct prompt")
	}
}

func commands() []cli.Command {
	return []cli.Command{
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
						execute("logs -f", psaIds, psaNames, "Which container would you like to see the logs of?")
					} else {
						execute("logs", psaIds, psaNames, "Which container would you like to see the logs of?")
					}
				} else {
					passThroughDocker()
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
					execute("start", psaIds, psaNames, "Which container would you like to start?")
				} else {
					passThroughDocker()
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
					execute("restart", psIds, psNames, "Which container would you like to restart?")
				} else {
					passThroughDocker()
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
					execute("stop", psIds, psNames, "Which container would you like to stop?")
				} else {
					passThroughDocker()
				}
				return nil
			},
		},
		{
			Name:  "ssh",
			Usage: "SSH into a container",
			Action: func(c *cli.Context) error {
				execute("ssh", psIds, psNames, "Which container would you like to connect with?")
				return nil
			},
		},
		{
			Name:  "env",
			Usage: "See the environment variables of a running container",
			Action: func(c *cli.Context) error {
				execute("env", psIds, psNames, "Which container would you like to see the environment variables of?")
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
					execute("rm", psaIds, psaNames, "Which container would you like to remove?")
				} else {
					passThroughDocker()
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
					execute("rmi", imageIds, imageNames, "Which image would you like to remove?")
				} else {
					passThroughDocker()
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
					execute("history", imageIds, imageNames, "Which image would you like to see the history of?")
				} else {
					passThroughDocker()
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
					Usage: "Select the container from a select docker",
				},
			},
			Action: func(c *cli.Context) error {
				if len(c.Args()) == 0 && c.Bool("s") {
					execute("stats", psIds, psNames, "Which container would you like to see that stats of?")
				} else {
					passThroughDocker()
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
					execute("inspect", psIds, psNames, "Which container would you like to inspect?")
				} else {
					passThroughDocker()
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
}
