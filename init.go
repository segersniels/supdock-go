package main

import (
	"os"
	"strings"

	util "github.com/segersniels/goutil"
)

func init() {
	commandNames := extractNames(commands())
	utilNames := []string{"-h", "--help", "-v", "--version"}
	if len(os.Args) > 1 && util.Exists(commandNames, os.Args[1]) && !util.Exists(utilNames, os.Args[1]) {
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
}
