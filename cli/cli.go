package cli

import (
	"os"

	util "github.com/segersniels/goutil"
	"github.com/segersniels/supdock/command"
	"github.com/urfave/cli"
)

type CLI struct {
	Commands []cli.Command
}

func (c CLI) Run() {
	app := cli.NewApp()
	app.Name = "supdock"
	app.Usage = "What's Up Dock(er)?"
	app.Version = "0.1.5"
	app.Commands = c.Commands

	names := command.ExtractNames(c.Commands)

	if len(os.Args) > 1 && util.Exists(names, os.Args[1]) {
		err := app.Run(os.Args)
		if err != nil {
			util.Warn(err)
			os.Exit(0)
		}
	} else {
		command.PassThrough()
	}
}
