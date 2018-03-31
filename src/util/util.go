package util

import (
	"fmt"
	"log"
	"os/exec"
	"reflect"
)

func usage() {
	output := `Usage: supdock [options] [command]

Options:
	-h, --help         output usage information

Commands:
	stop              Stop a running container
	start             Start a stopped container
	logs              See the logs of a container
	rm                Remove a container
	rmi               Remove an image
	prune             Remove stopped containers and dangling images
	stats             See the stats of a container
	ssh               SSH into a container
	history           See the history of an image
	env               See the environment variables of a running container
`
	fmt.Print(output)
}

// Help : call docker help as output for supdock
func Help() {
	usage()
	dockerOut, err := exec.Command("docker", "--help").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", dockerOut)
}

// SliceExists : check if a slice contains a specific value
func SliceExists(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("SliceExists() given a non-slice type")
	}
	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}
	return false
}
