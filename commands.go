package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func passThroughDocker() {
	cmd := exec.Command("docker", os.Args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func customDocker(args []string) {
	var errbuf bytes.Buffer
	if exists(args, "|") {
		index := getIndex(args, "|")
		length := len(args)
		cmd := exec.Command("docker", args[0:index]...)
		pipeCmd := exec.Command(args[index+1], args[index+2:length]...)

		pipeCmd.Stdin, _ = cmd.StdoutPipe()
		pipeCmd.Stdout = os.Stdout
		pipeCmd.Stderr = &errbuf

		err := pipeCmd.Start()
		if err != nil {
			log.Fatal(strings.TrimSpace(errbuf.String()))
		}

		err = cmd.Run()
		if err != nil {
			log.Fatal(strings.TrimSpace(errbuf.String()))
		}

		err = pipeCmd.Wait()
		if err != nil {
			log.Fatal(strings.TrimSpace(errbuf.String()))
		}
	} else {
		cmd := exec.Command("docker", args...)
		cmd.Stderr = &errbuf
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		err := cmd.Run()
		if err != nil {
			log.Fatal(strings.TrimSpace(errbuf.String()))
		}
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
		choice := id + " - " + strings.TrimLeft(names[index], "/")
		if choice != " - " {
			choices = append(choices, choice)
		}
	}
	return choices
}

func selectID(ids []string, names []string, question string) string {
	options := constructChoices(ids, names)
	answer := promptQuestion(question, options)
	id := strings.Split(answer, " - ")[0]
	return id
}

func execute(command string, ids []string, names []string, question string) {
	if len(ids) >= 1 && len(names) >= 1 {
		id := selectID(ids, names, question)
		switch command {
		case "ssh":
			shell := promptQuestion("Which shell is the container using?", []string{"bash", "ash"})
			customDocker([]string{"exec", "-ti", id, shell})
		case "env":
			customDocker([]string{"exec", "-ti", id, "env"})
		case "logs -f":
			customDocker([]string{"logs", "-f", id})
		case "stats --no-stream":
			customDocker([]string{"stats", "--no-stream", id})
		default:
			customDocker([]string{command, id})
		}
	} else {
		log.Fatal("No options found to construct prompt")
	}
}

func start(id string) {
	err := docker.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}

func stop(id string) {
	err := docker.ContainerStop(context.Background(), id, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}

func restart(id string) {
	err := docker.ContainerRestart(context.Background(), id, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}

func remove(removeType string, id string) {
	var err error
	switch removeType {
	case "container":
		err = docker.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{})
	case "image":
		_, err = docker.ImageRemove(context.Background(), id, types.ImageRemoveOptions{})
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
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
					id := selectID(psaIds, psaNames, "Which container would you like to start?")
					start(id)
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
					id := selectID(psIds, psNames, "Which container would you like to restart?")
					restart(id)
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
					if psIds == nil {
						log.Fatal("No containers found to start")
					}
					id := selectID(psIds, psNames, "Which container would you like to stop?")
					stop(id)
				} else {
					passThroughDocker()
				}
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:  "all",
					Usage: "Stop all running containers",
					Action: func(c *cli.Context) error {
						if psIds == nil {
							log.Fatal("No containers found to start")
						}
						for _, id := range psIds {
							stop(id)
						}
						return nil
					},
				},
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
					id := selectID(psaIds, psaNames, "Which container would you like to remove?")
					remove("container", id)
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
					id := selectID(imageIds, imageNames, "Which image would you like to remove?")
					remove("image", id)
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
				if c.Bool("s") {
					if c.Bool("no-stream") {
						execute("stats --no-stream", psIds, psNames, "Which container would you like to see the stats of?")
					} else {
						execute("stats", psIds, psNames, "Which container would you like to see the stats of?")
					}
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
				customDocker([]string{"system", "prune", "-f"})
				return nil
			},
		},
	}
}
