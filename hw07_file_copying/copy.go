package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func copyImpl(fIn io.Reader, fOut io.Writer, limit int) error {
	pb := NewProgressBar("copying...")

	nCopy := 0
	bufSize := 1024
	if bufSize > limit {
		bufSize = limit
	}

	pb.Start(int64(limit))

	eof := false
	for nCopy < limit && !eof {
		buf := make([]byte, bufSize)
		nRead, err := io.ReadAtLeast(fIn, buf, bufSize)
		if errors.Is(err, io.ErrUnexpectedEOF) {
			eof = true
		} else if err != nil {
			return err
		}

		nWrite := 0
		for nWrite < nRead {
			nOut, err := fOut.Write(buf[nWrite:nRead])
			if err != nil {
				return err
			}
			nWrite += nOut
		}

		nCopy += nRead
		pb.Update(int64(nCopy))
	}

	pb.Finish()
	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	// open input file
	fIn, err := os.OpenFile(fromPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer fIn.Close()

	// check is input file not a directory
	fi, err := fIn.Stat()
	if err != nil || fi.IsDir() {
		return ErrUnsupportedFile
	}

	// check offset and limit
	if fi.Size() < offset {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 {
		limit = fi.Size() - offset
	}

	// open or create out file
	fOut, err := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if errors.Is(err, syscall.EISDIR) {
		// create the file with the same name
		fOutName := filepath.Join(toPath, filepath.Base(toPath))
		fOut, err = os.OpenFile(fOutName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	}
	if err != nil {
		return err
	}
	defer fOut.Close()

	// do not copy if size == 0
	if fi.Size() == 0 {
		return nil
	}

	if _, err = fIn.Seek(offset, 0); err != nil {
		return err
	}

	return copyImpl(fIn, fOut, int(limit))
}
