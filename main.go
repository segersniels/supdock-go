package main

import (
	"github.com/segersniels/supdock/cli"
	"github.com/segersniels/supdock/command"
)

func main() {
	cli := cli.CLI{
		Commands: command.Commands(),
	}
	cli.Run()
}
