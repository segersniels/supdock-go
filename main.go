package main

import (
	"os"

	util "github.com/segersniels/goutil"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "supdock"
	app.Usage = "What's Up Dock(er)?"
	app.Version = "0.1.5"
	app.Commands = commands()

	names := extractNames(app.Commands)
	utils := []string{"-h", "--help", "-v", "--version"}

	if len(os.Args) > 1 && (util.Exists(names, os.Args[1]) || util.Exists(utils, os.Args[1])) {
		err := app.Run(os.Args)
		if err != nil {
			util.Warn(err)
			os.Exit(0)
		}
	} else {
		passThrough()
	}
}
