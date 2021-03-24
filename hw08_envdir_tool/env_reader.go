package main

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrNotDirectory = errors.New("path is not a directory")

func findEndIdx(buf []byte) int {
	if len(buf) == 0 {
		return 0
	}

	endIdx := len(buf)
	for endIdx > 0 {
		if buf[endIdx-1] == byte(' ') || buf[endIdx-1] == byte('\t') {
			endIdx--
		} else {
			break
		}
	}

	// find terminated nulls
	idx := 0
	for idx < endIdx {
		if buf[idx] == 0x00 {
			buf[idx] = '\n'
		}
		idx++
	}

	return endIdx
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirI, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !dirI.IsDir() {
		return nil, ErrNotDirectory
	}

	env := Environment{}
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Warning(err)
			return nil
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.OpenFile(path, os.O_RDONLY, os.ModeExclusive)
		if err != nil {
			log.WithField("path", path).Warning(err)
			return nil
		}
		defer f.Close()

		bF := bufio.NewReader(f)
		buf := make([]byte, 0)
		readMore := true
		for readMore {
			var l []byte
			l, readMore, err = bF.ReadLine()
			if err != nil {
				log.WithField("path", path).Warning("failed to read line: ", err)
			} else {
				buf = append(buf, l...)
			}
		}

		endIdx := findEndIdx(buf)

		fName := filepath.Base(path)
		env[fName] = EnvValue{string(buf[:endIdx]), endIdx == 0}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return env, nil
}
