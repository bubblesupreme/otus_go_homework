package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func compareFiles(t *testing.T, f1, f2 string) {
	var outbuf, errbuf bytes.Buffer
	cmp := exec.Cmd{Path: "/bin/cmp", Args: []string{"", f1, f2}}
	cmp.Stdout = &outbuf
	cmp.Stderr = &errbuf
	fmt.Println(cmp.String())
	err := cmp.Run()
	require.Nil(t, err, "\nerr:%s\nstderr:%s\nstdout:%s", err, errbuf, outbuf)
}

func TestCopyNonErr(t *testing.T) {
	input := "testdata/input.txt"
	empty := "testdata/empty.txt"
	tests := []struct {
		name   string
		from   string
		offset int64
		limit  int64
		cmp    string
	}{
		{name: "offset0 limit0", from: input, offset: 0, limit: 0, cmp: "testdata/out_offset0_limit0.txt"},
		{name: "offset0 limit10", from: input, offset: 0, limit: 10, cmp: "testdata/out_offset0_limit10.txt"},
		{name: "offset0 limit1000", from: input, offset: 0, limit: 1000, cmp: "testdata/out_offset0_limit1000.txt"},
		{name: "offset0 limit10000", from: input, offset: 0, limit: 10000, cmp: "testdata/out_offset0_limit10000.txt"},
		{name: "offset100 limit1000", from: input, offset: 100, limit: 1000, cmp: "testdata/out_offset100_limit1000.txt"},
		{name: "offset6000 limit1000", from: input, offset: 6000, limit: 1000, cmp: "testdata/out_offset6000_limit1000.txt"},
		{name: "empty offset0 limit0", from: empty, offset: 0, limit: 0, cmp: "testdata/out_empty.txt"},
		{name: "empty offset0 limit10", from: empty, offset: 0, limit: 10, cmp: "testdata/out_empty.txt"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			out, err := ioutil.TempFile("testdata/", "")
			if err != nil {
				log.Fatal("Can't create the output file: ", err)
			}
			defer os.Remove(out.Name())

			err = Copy(tc.from, out.Name(), tc.offset, tc.limit)
			require.Nil(t, err)

			compareFiles(t, out.Name(), tc.cmp)
		})
	}
}

func TestCopyDirectory(t *testing.T) {
	tests := []struct {
		name   string
		from   string
		to     string
		offset int64
		limit  int64
		err    error
	}{
		{name: "directory as input", from: "testdata/", to: "", offset: 0, limit: 10, err: ErrUnsupportedFile},
		{name: "directory as out", from: "testdata/input.txt", to: "testdata/", offset: 0, limit: 1000, err: nil},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := Copy(tc.from, tc.to, tc.offset, tc.limit)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestCopyErrExpected(t *testing.T) {
	tests := []struct {
		name   string
		from   string
		offset int64
		limit  int64
		err    error
	}{
		{name: "offset exceeds file size", from: "testdata/empty.txt", offset: 1, limit: 0, err: ErrOffsetExceedsFileSize},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			out, err := ioutil.TempFile("testdata/", "")
			if err != nil {
				log.Fatal("Can't create the output file: ", err)
			}
			defer os.Remove(out.Name())

			err = Copy(tc.from, out.Name(), tc.offset, tc.limit)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestCopyPermission(t *testing.T) {
	dir, err := ioutil.TempDir("testdata/", "")
	if err != nil {
		log.Fatal("Can't create the output file: ", err)
	}
	defer os.RemoveAll(dir)

	t.Run("don't write permission", func(t *testing.T) {
		fName := filepath.Join(dir, "read_permission.txt")
		f, err := os.OpenFile(fName, os.O_RDONLY|os.O_CREATE, 0400)
		if err != nil {
			log.Fatal("Can't create the output file: ", err)
		}
		defer f.Close()

		err = Copy("testdata/input.txt", f.Name(), 1, 1)
		require.ErrorIs(t, err, syscall.EACCES)
	})

	t.Run("don't read permission", func(t *testing.T) {
		fName := filepath.Join(dir, "write_permission.txt")
		f, err := os.OpenFile(fName, os.O_WRONLY|os.O_CREATE, 0200)
		if err != nil {
			log.Fatal("Can't create the output file: ", err)
		}
		defer f.Close()

		err = Copy(f.Name(), "", 1, 1)
		require.ErrorIs(t, err, syscall.EACCES)
	})
}

func TestCopyCreate(t *testing.T) {
	tests := []struct {
		name   string
		from   string
		offset int64
		limit  int64
		err    error
	}{
		{name: "output file doesn't exist", from: "testdata/empty.txt", offset: 1, limit: 0, err: ErrOffsetExceedsFileSize},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("testdata/", "")
			if err != nil {
				log.Fatal("Can't create the output directory: ", err)
			}
			defer os.RemoveAll(dir)

			fName := filepath.Join(dir, "out.txt")
			err = Copy(tc.from, fName, tc.offset, tc.limit)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestCopyRewrite(t *testing.T) {
	tests := []struct {
		name    string
		from    string
		offset1 int64
		offset2 int64
		limit   int64
		cmp1    string
		cmp2    string
	}{
		{name: "rewrite output file", from: "testdata/input.txt", offset1: 0, offset2: 100, limit: 1000, cmp1: "testdata/out_offset0_limit1000.txt", cmp2: "testdata/out_offset100_limit1000.txt"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			out, err := ioutil.TempFile("testdata/", "")
			if err != nil {
				log.Fatal("Can't create the output file: ", err)
			}
			defer os.Remove(out.Name())

			err = Copy(tc.from, out.Name(), tc.offset1, tc.limit)
			require.Nil(t, err)
			compareFiles(t, out.Name(), tc.cmp1)

			err = Copy(tc.from, out.Name(), tc.offset2, tc.limit)
			require.Nil(t, err)
			compareFiles(t, out.Name(), tc.cmp2)
		})
	}
}
