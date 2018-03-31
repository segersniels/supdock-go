package util

import (
	"fmt"
	"log"
	"os/exec"
	"reflect"
)

// Help : call docker help as output for supdock
func Help() {
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
