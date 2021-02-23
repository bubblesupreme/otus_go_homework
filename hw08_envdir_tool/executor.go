//nolint:errorlint
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

var errParseEnvVar = errors.New("don't environment variable")

// Function parseEnvVar get environment variable in the form "key=value"
// return key, value, error.
func parseEnvVar(v string) (string, string, error) {
	res := strings.SplitN(v, "=", 2)
	if len(res) != 2 {
		return "", "", errParseEnvVar
	}
	return res[0], res[1], nil
}

func createEnvVar(k, v string) string {
	return fmt.Sprintf("%s=%s", k, v)
}

// Function getOsEnv get current OS environment and return it in the form map[key]value.
func getOsEnv() map[string]string {
	osEnv := os.Environ()

	res := make(map[string]string, len(osEnv))
	for _, e := range osEnv {
		k, v, err := parseEnvVar(e)
		if err != nil {
			log.WithField("env var", e).Fatal("failed to parse environment variable: ", err)
		} else {
			res[k] = v
		}
	}

	return res
}

// Function getEnv get current environment and update by Environment struct
// return a strings representing the environment in the form "key=value".
func getEnv(env Environment) []string {
	osEnv := getOsEnv()
	for k, newVal := range env {
		_, ok := osEnv[k]
		if newVal.NeedRemove && ok {
			delete(osEnv, k)
		} else if !newVal.NeedRemove {
			osEnv[k] = newVal.Value
		}
	}

	res := make([]string, len(osEnv))
	idx := 0
	for k, v := range osEnv {
		res[idx] = createEnvVar(k, v)
		idx++
	}

	return res
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	updatedEnv := getEnv(env)

	exCmd := exec.Cmd{Path: cmd[0], Args: cmd, Stderr: os.Stderr, Stdout: os.Stdout, Env: updatedEnv}
	if err := exCmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
			log.WithField("error", exiterr).Fatal("failed to get the wait status")
		}
		log.WithField("error", err).Fatal("failed to cast to ExitError")
	}

	return 0
}
