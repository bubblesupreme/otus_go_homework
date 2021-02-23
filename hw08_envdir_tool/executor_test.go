//nolint:ifshort
package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

const defaultEnvVal = "test"

func createEchoScript(vars []string) (*os.File, error) {
	str := "#!/usr/bin/env bash\necho -e \""
	for _, k := range vars {
		str += fmt.Sprintf("%s is ${%s}\n", k, k)
	}
	str += "\""

	f, err := ioutil.TempFile("testdata/", "test.*.sh")
	if err != nil {
		return nil, err
	}

	if err := f.Chmod(0500); err != nil {
		return nil, err
	}

	nWr := 0
	for nWr < len(str) {
		n, err := f.WriteString(str)
		if err != nil {
			return nil, err
		}
		nWr += n
	}

	return f, nil
}

func setEnvVar(name string) {
	if err := os.Setenv(name, defaultEnvVal); err != nil {
		log.Fatal("failed to set environment variable: ", name)
	}
}

func unsetEnvVar(name string) {
	if err := os.Unsetenv(name); err != nil {
		log.Fatal("failed to unset environment variable: ", name)
	}
}

func TestRunCmd(t *testing.T) {
	name1 := "full test"
	env1 := Environment{
		"BAR":           EnvValue{"bar", false},
		"BAR_DEL":       EnvValue{"bar", true},
		"EMPTY":         EnvValue{"", false},
		"EMPTY_DEL":     EnvValue{"", true},
		"EQUAL_SMBL":    EnvValue{"equal=equal", false},
		"NOT_EXIST":     EnvValue{"not exist", false},
		"NOT_EXIST_REM": EnvValue{"not exist remove", true},
	}
	// order is important
	expected1 := "BAR is bar\nBAR_DEL is \nEMPTY is \nEMPTY_DEL is \nEQUAL_SMBL is equal=equal\nNOT_EXIST is not exist\nNOT_EXIST_REM is \n\n"

	keys1 := make([]string, 0, len(env1))
	for k := range env1 {
		keys1 = append(keys1, k)
	}
	// sort for correct order
	sort.Strings(keys1)

	name2 := "empty env"
	env2 := Environment{}
	keys2 := []string{"BAR"}
	expected2 := fmt.Sprintf("BAR is %s\n\n", defaultEnvVal)

	tests := []struct {
		name     string
		env      Environment
		keys     []string
		expected string
	}{
		{name: name1, env: env1, keys: keys1, expected: expected1},
		{name: name2, env: env2, keys: keys2, expected: expected2},
	}

	setEnvVar("BAR")
	setEnvVar("BAR_DEL")
	setEnvVar("EMPTY")
	unsetEnvVar("NOT_EXIST")

	shell := os.Getenv("SHELL")
	if shell == "" {
		log.Fatal("failed to find shell")
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			f, err := createEchoScript(tc.keys)
			if err != nil {
				log.Fatal("failed to create temporary script file")
			}
			defer os.Remove(f.Name())

			cmd := make([]string, 2)
			cmd[0] = shell
			cmd[1] = f.Name()

			realStdout := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				log.Fatal("failed to creat a pipe")
			}
			os.Stdout = w

			returnCode := RunCmd(cmd, tc.env)
			if err := w.Close(); err != nil {
				log.Fatal("failed to close a pipe")
			}

			out, err := ioutil.ReadAll(r)
			if err := r.Close(); err != nil {
				log.Fatal("failed to close a pipe")
			}

			if err != nil && errors.Is(err, io.EOF) {
				log.Fatal("failed to read output")
			}
			os.Stdout = realStdout

			require.Equal(t, 0, returnCode)
			require.Equal(t, tc.expected, string(out))
		})
	}
}
