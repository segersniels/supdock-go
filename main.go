package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "supdock"
	app.Usage = "What's Up Dock(er)?"
	app.Version = "0.1.7-rc.2"
	app.Commands = commands()

	commandNames := extractNames(app.Commands)
	utilNames := []string{"-h", "--help", "-v", "--version"}

	if len(os.Args) > 1 && (exists(commandNames, os.Args[1]) || exists(utilNames, os.Args[1])) {
		err := app.Run(os.Args)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		passThroughDocker()
	}
}
