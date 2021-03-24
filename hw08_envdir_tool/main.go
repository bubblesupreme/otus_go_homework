package main

import (
	"errors"
	"log"
	"os"
)

var ErrInvalidCmdArgs = errors.New("invalid command arguments")

func main() {
	if len(os.Args) < 3 {
		log.Fatal(ErrInvalidCmdArgs)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	returnCode := RunCmd(os.Args[2:], env)
	os.Exit(returnCode)
}
