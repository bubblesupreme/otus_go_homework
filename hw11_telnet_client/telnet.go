package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"syscall"
	"time"
)

const (
	bufSize           = 1024 * 4
	network           = "tcp"
	defaultTimeoutStr = "10s"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClientImpl{address, in, out, timeout, nil, nil, nil, true}
}

func NewTelnetClientWithContext(ctx context.Context,
	cancel context.CancelFunc,
	address string,
	timeout time.Duration,
	in io.ReadCloser,
	out io.Writer) TelnetClient {
	return &telnetClientImpl{address, in, out, timeout, nil, ctx, cancel, false}
}

type telnetClientImpl struct {
	address      string
	in           io.ReadCloser
	out          io.Writer
	timeout      time.Duration
	connection   net.Conn
	ctx          context.Context
	cancel       context.CancelFunc
	isCtxDefault bool
}

func (c *telnetClientImpl) Connect() error {
	var err error
	c.connection, err = net.DialTimeout(network, c.address, c.timeout)
	if c.isCtxDefault {
		c.ctx, c.cancel = context.WithCancel(context.Background())
	}

	if err == nil {
		go c.callTelnet()
	}

	return err
}

func (c *telnetClientImpl) Send() error {
	readMore := true
	for readMore {
		buf := make([]byte, bufSize)
		nRead, err := c.in.Read(buf)
		if nRead == 0 || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
			readMore = false
		} else if err != nil {
			return err
		}

		nWrite := 0
		for nWrite < nRead {
			nOut, err := c.connection.Write(buf[nWrite:nRead])
			if err != nil {
				return err
			}
			nWrite += nOut
		}
	}

	return nil
}

func (c *telnetClientImpl) Receive() error {
	buf := make([]byte, bufSize)
	nRead, err := c.connection.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	nWrite := 0
	for nWrite < nRead {
		nOut, err := c.out.Write(buf[nWrite:nRead])
		if err != nil {
			return err
		}
		nWrite += nOut
	}

	return nil
}

func (c *telnetClientImpl) Close() error {
	c.cancel()
	return c.connection.Close()
}

func (c *telnetClientImpl) callTelnet() {
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				err := c.Receive()
				if errors.Is(err, syscall.ECONNRESET) {
					log.Printf("attempt to read from a closed server")
					_ = c.Close()
					return
				} else if err != nil {
					log.Printf("failed to receive message: %s", err)
				}
			}
		}
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			err := c.Send()
			if errors.Is(err, syscall.EPIPE) {
				log.Printf("attempt to write to a closed server")
				_ = c.Close()
				return
			} else if err != nil {
				log.Printf("failed to send message: %s", err)
			}
		}
	}
}
