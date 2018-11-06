package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

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
						executeDocker("logs-force", psaIds, psaNames, "Which container would you like to see the logs of?")
					} else {
						executeDocker("logs", psaIds, psaNames, "Which container would you like to see the logs of?")
					}
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
				executeDocker("ssh", psIds, psNames, "Which container would you like to connect with?")
				return nil
			},
		},
		{
			Name:  "env",
			Usage: "See the environment variables of a running container",
			Action: func(c *cli.Context) error {
				executeDocker("env", psIds, psNames, "Which container would you like to see the environment variables of?")
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
				if len(c.Args()) == 0 {
					id := selectID(imageIds, imageNames, "Which image would you like to remove?")
					if c.NumFlags() == 0 {
						remove("image", id)
					} else if c.Bool("f") {
						remove("image-force", id)
					}
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
					executeDocker("history", imageIds, imageNames, "Which image would you like to see the history of?")
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
						executeDocker("stats-no-stream", psIds, psNames, "Which container would you like to see the stats of?")
					} else {
						executeDocker("stats", psIds, psNames, "Which container would you like to see the stats of?")
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
					executeDocker("inspect", psIds, psNames, "Which container would you like to inspect?")
				} else {
					passThroughDocker()
				}
				return nil
			},
		},
		{
			Name:  "prune",
			Usage: "Remove stopped containers and dangling images. For more detailed usage refer to 'docker system prune -h'",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "a, all",
					Usage: "Remove all unused images not just dangling ones",
				},
				cli.BoolFlag{
					Name:  "f, force",
					Usage: "Do not prompt for confirmation",
				},
			},
			Action: func(c *cli.Context) error {
				var flags []string
				for _, flag := range c.FlagNames() {
					if c.IsSet(flag) {
						flags = append(flags, "-"+flag)
					}
				}
				customDocker(append([]string{"system", "prune"}, flags...))
				return nil
			},
		},
		{
			Name:  "compose",
			Usage: "Allows for dynamic docker-compose usage",
			Subcommands: []cli.Command{
				{
					Name:  "build",
					Usage: "Build or rebuild services",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:        "depth",
							Usage:       "Define the depth of the docker-compose file search",
							Value:       6,
							Destination: &depth,
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 && c.NumFlags() == 0 {
							executeCompose("build", "Which project would you like to build?")
						} else {
							passThroughCompose()
						}
						return nil
					},
				},
				{
					Name:  "restart",
					Usage: "Restart services",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:        "depth",
							Usage:       "Define the depth of the docker-compose file search",
							Value:       4,
							Destination: &depth,
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 && c.NumFlags() == 0 {
							executeCompose("restart", "Which project would you like to restart?")
						} else {
							passThroughCompose()
						}
						return nil
					},
				},
				{
					Name:  "pull",
					Usage: "Pull service images",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:        "depth",
							Usage:       "Define the depth of the docker-compose file search",
							Value:       4,
							Destination: &depth,
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 && c.NumFlags() == 0 {
							executeCompose("pull", "Which project would you like to pull the images from?")
						} else {
							passThroughCompose()
						}
						return nil
					},
				},
				{
					Name:  "start",
					Usage: "Start services",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:        "depth",
							Usage:       "Define the depth of the docker-compose file search",
							Value:       4,
							Destination: &depth,
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 && c.NumFlags() == 0 {
							executeCompose("start", "Which project would you like to start?")
						} else {
							passThroughCompose()
						}
						return nil
					},
				},
				{
					Name:  "up",
					Usage: "Create and start containers",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "d, detached",
							Usage: "Detached mode: Run containers in the background",
						},
						cli.IntFlag{
							Name:        "depth",
							Usage:       "Define the depth of the docker-compose file search",
							Value:       4,
							Destination: &depth,
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 {
							if c.Bool("d") {
								executeCompose("up-detached", "Which project would you like to start?")
							} else {
								executeCompose("up", "Which project would you like to start?")
							}
						} else {
							passThroughCompose()
						}
						return nil
					},
				},
				{
					Name:  "down",
					Usage: "Stop and remove containers, networks, images, and volumes",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:        "depth",
							Usage:       "Define the depth of the docker-compose file search",
							Value:       4,
							Destination: &depth,
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 && c.NumFlags() == 0 {
							executeCompose("down", "Which project would you like to bring down?")
						} else {
							passThroughCompose()
						}
						return nil
					},
				},
				{
					Name:  "stop",
					Usage: "Stop services",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:        "depth",
							Usage:       "Define the depth of the docker-compose file search",
							Value:       4,
							Destination: &depth,
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 && c.NumFlags() == 0 {
							executeCompose("stop", "Which project would you like to stop?")
						} else {
							passThroughCompose()
						}
						return nil
					},
				},
				{
					Name:  "top",
					Usage: "Display the running processes",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:        "depth",
							Usage:       "Define the depth of the docker-compose file search",
							Value:       4,
							Destination: &depth,
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 && c.NumFlags() == 0 {
							executeCompose("top", "Which project would you like to see the running processes of?")
						} else {
							passThroughCompose()
						}
						return nil
					},
				},
				{
					Name:  "logs",
					Usage: "View output from containers",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:        "depth",
							Usage:       "Define the depth of the docker-compose file search",
							Value:       4,
							Destination: &depth,
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 && c.NumFlags() == 0 {
							executeCompose("logs", "Which project would you like to see the logs of?")
						} else {
							passThroughCompose()
						}
						return nil
					},
				},
				{
					Name:  "ps",
					Usage: "List containers",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:        "depth",
							Usage:       "Define the depth of the docker-compose file search",
							Value:       4,
							Destination: &depth,
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 && c.NumFlags() == 0 {
							executeCompose("ps", "Which project would you like to see the running processes of?")
						} else {
							passThroughCompose()
						}
						return nil
					},
				},
				{
					Name:    "list",
					Aliases: []string{"ls"},
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:        "depth",
							Usage:       "Define the depth of the docker-compose file search",
							Value:       4,
							Destination: &depth,
						},
					},
					Usage: "List all your docker-compose projects",
					Action: func(c *cli.Context) error {
						files := searchComposeFiles()
						projects, _ := yaml.Marshal(files)
						fmt.Println(string(projects))
						return nil
					},
				},
			},
		},
	}
}
