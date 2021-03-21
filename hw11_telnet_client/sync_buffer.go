package main

import (
	"bytes"
	"sync"
)

type Buffer struct {
	buf bytes.Buffer
	m   sync.Mutex
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.buf.Read(p)
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.buf.Write(p)
}

func (b *Buffer) String() string {
	b.m.Lock()
	defer b.m.Unlock()
	return b.buf.String()
}
