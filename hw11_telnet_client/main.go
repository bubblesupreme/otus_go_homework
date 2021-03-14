package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatalf("incompatible command arguments: expected %d, but received %d", 2, len(args))
	}

	timeoutStr := flag.String("timeout", defaultTimeoutStr, "specifier connection timeout")
	flag.Parse()

	timeout, err := time.ParseDuration(*timeoutStr)
	if err != nil {
		log.Fatalf("failed to parse timeout %s", *timeoutStr)
	}

	address := flag.Arg(0) + ":" + flag.Arg(1)

	ctx, cancel := context.WithCancel(context.Background())
	client := NewTelnetClientWithContext(ctx, cancel, address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatalf("failed to connect to %s", address)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			select {
			case s, ok := <-sigs:
				if !ok {
					return
				}

				if s == syscall.SIGINT {
					if err := client.Close(); err != nil {
						log.Println("failed to close connection")
					}

					log.Println("...Connection was closed by sigint")
					return
				}
			default:
				break
			}
		}
	}
}
