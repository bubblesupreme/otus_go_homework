package main

import (
	"errors"
	"log"
	"os"
)

var ErrInvalidCmdArgs = errors.New("invalid command arguments")

func main() {
	if len(os.Args) < 2 {
		log.Fatal(ErrInvalidCmdArgs)
	}

	_, err := ReadDir(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}
}
