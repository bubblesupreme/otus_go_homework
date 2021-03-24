package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDirFull(t *testing.T) {
	inputDir := "testdata/env/"
	expected := Environment{
		"BAR":        EnvValue{"bar", false},
		"EMPTY":      EnvValue{"", true},
		"FOO":        EnvValue{"   foo\nwith new line", false},
		"HELLO":      EnvValue{"\"hello\"", false},
		"SPACES":     EnvValue{"", true},
		"TABULATION": EnvValue{"", true},
		"ONE_LINE":   EnvValue{"\"one line\"", false},
		"UNSET":      EnvValue{"", true},
	}
	size := len(expected)

	env, err := ReadDir(inputDir)
	require.NoError(t, err)
	require.Equal(t, len(env), size)
	for v, r := range env {
		exp, ok := expected[v]
		require.True(t, ok)
		if ok {
			require.Equal(t, exp, r)
		}
	}
}

func TestReadDirEmpty(t *testing.T) {
	inputDir, err := ioutil.TempDir("testdata/", "")
	if err != nil {
		log.Fatalf("Can't create a directory")
	}
	defer os.RemoveAll(inputDir)

	size := 0

	env, err := ReadDir(inputDir)
	require.NoError(t, err)
	require.Equal(t, len(env), size)
}

func TestReadDirNotDirectory(t *testing.T) {
	inputDir := "testdata/env/BAR"

	env, err := ReadDir(inputDir)
	require.ErrorIs(t, err, ErrNotDirectory)
	require.Nil(t, env)
}

func TestReadDirNotExist(t *testing.T) {
	inputDir := "./foo"

	env, err := ReadDir(inputDir)
	require.Error(t, err)
	require.Nil(t, env)
}
