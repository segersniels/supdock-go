package main

import (
	"context"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

var psIds, psaIds, imageIds, psNames, psaNames, imageNames []string
var docker *client.Client
var depth int

func getContainerInformation(cli *client.Client, all bool) ([]string, []string) {
	var ids, names []string
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: all})
	if err != nil {
		log.Fatal(err)
	}
	for _, container := range containers {
		ids = append(ids, container.ID[0:12])
		names = append(names, container.Names...)
	}
	return ids, names
}

func getImageInformation(cli *client.Client) ([]string, []string) {
	var ids, names []string
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, image := range images {
		ids = append(ids, strings.SplitAfter(image.ID, ":")[1][0:12])
		names = append(names, image.RepoTags...)
	}
	return ids, names
}

func init() {
	commandNames := extractNames(commands())
	utilNames := []string{"-h", "--help", "-v", "--version"}
	if len(os.Args) > 1 && exists(commandNames, os.Args[1]) && !exists(utilNames, os.Args[1]) {
		docker, _ = client.NewEnvClient()
		psIds, psNames = getContainerInformation(docker, false)
		psaIds, psaNames = getContainerInformation(docker, true)
		imageIds, imageNames = getImageInformation(docker)
	}
}
