package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClientBasic(t *testing.T) {
	str1 := "hello\n"
	str2 := "world\n"

	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer func() { require.NoError(t, l.Close()) }()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("10s")
		require.NoError(t, err)

		client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
		require.NoError(t, client.Connect())
		defer func() { require.NoError(t, client.Close()) }()
		in.WriteString(str1)
		err = client.Send()
		require.NoError(t, err)

		err = client.Receive()
		require.NoError(t, err)
		require.Equal(t, str2, out.String())
	}()

	go func() {
		defer wg.Done()

		conn, err := l.Accept()
		require.NoError(t, err)
		require.NotNil(t, conn)
		defer func() { require.NoError(t, conn.Close()) }()

		request := make([]byte, 1024)
		n, err := conn.Read(request)
		require.NoError(t, err)
		require.Equal(t, str1, string(request)[:n])

		n, err = conn.Write([]byte(str2))
		require.NoError(t, err)
		require.NotEqual(t, 0, n)
	}()

	wg.Wait()
}

func TestTelnetClientNoExistServer(t *testing.T) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	timeout, err := time.ParseDuration("10s")
	require.NoError(t, err)
	client := NewTelnetClient("noexist:10", timeout, ioutil.NopCloser(in), out)
	require.EqualError(t, client.Connect(), "dial tcp: lookup noexist: no such host")
}

func TestTelnetClientClosedClient(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer func() { require.NoError(t, l.Close()) }()

	in := &bytes.Buffer{}
	out := &bytes.Buffer{}

	timeout, err := time.ParseDuration("10s")
	require.NoError(t, err)

	client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
	require.NoError(t, client.Connect())

	require.NoError(t, client.Close())
	time.Sleep(3 * time.Second)
	in.WriteString("hello\n")
	err = client.Send()
	require.Error(t, err)
}

func TestTelnetClientClosedServer(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)

	in := &bytes.Buffer{}
	out := &bytes.Buffer{}

	var buf bytes.Buffer
	log.SetOutput(&buf)

	timeout, err := time.ParseDuration("10s")
	require.NoError(t, err)

	client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
	require.NoError(t, client.Connect())
	defer func() { require.Error(t, client.Close()) }()

	require.NoError(t, l.Close())
	time.Sleep(3 * time.Second)

	in.WriteString("hello\n")

	res := strings.Contains(buf.String(), "attempt to write to a closed server") ||
		strings.Contains(buf.String(), "attempt to read from a closed server")
	require.True(t, res)
	t.Log(buf.String())
	log.SetOutput(os.Stderr)
}

func TestTelnetClientWithContext(t *testing.T) {
	const sleepSec = 3
	str1 := []byte("hello1\n")
	str2 := []byte("\n")
	checkStr := []byte("hello2\n")

	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer func() { require.NoError(t, l.Close()) }()

	timeout, err := time.ParseDuration("10s")
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	require.NoError(t, err)
	var conn net.Conn
	var ctx context.Context
	var cancel context.CancelFunc
	var client TelnetClient

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		ctx, cancel = context.WithCancel(context.Background())
		client = NewTelnetClientWithContext(ctx, cancel, l.Addr().String(), timeout, ioutil.NopCloser(in), out)
		require.NoError(t, client.Connect())
	}()

	go func() {
		defer wg.Done()
		conn, err = l.Accept()
		require.NoError(t, err)
		require.NotNil(t, conn)
	}()

	wg.Wait()

	defer func() { require.NoError(t, conn.Close()) }()
	defer func() { require.NoError(t, client.Close()) }()

	n, err := conn.Write(str1)
	require.NoError(t, err)
	require.NotEqual(t, 0, n)

	time.Sleep(sleepSec * time.Second)
	cancel()

	n, err = conn.Write(str2)
	require.NoError(t, err)
	require.NotEqual(t, 0, n)

	n, err = conn.Write(checkStr)
	require.NoError(t, err)
	require.NotEqual(t, 0, n)

	require.False(t, strings.Contains(out.String(), string(checkStr)))
}
