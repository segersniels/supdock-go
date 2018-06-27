package command

import (
	"fmt"
	"runtime"
	"strings"

	util "github.com/segersniels/goutil"
)

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
