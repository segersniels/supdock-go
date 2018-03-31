package util

import (
	"fmt"
	"log"
	"os/exec"
	"reflect"
)

func Help() {
	dockerOut, err := exec.Command("docker", "-h").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", dockerOut)
}

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
