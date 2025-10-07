package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var (
	timeout time.Duration
)

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		log.Fatalf("Usage: go-telnet [--timeout=10s] host port")
	}

	host := args[0]
	port := args[1]

	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatalf("connect error: %v", err)
	}
	defer func() {
		client.Close()
		fmt.Fprintln(os.Stderr, "...EOF")
	}()

	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	go func() {
		defer cancel()
		for {
			if err := client.Send(); err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Fprintln(os.Stderr, "...EOF")
					return
				}
				fmt.Fprintf(os.Stderr, "Send error: %v\n", err)
				return
			}
		}
	}()

	go func() {
		defer cancel()
		for {
			if err := client.Receive(); err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
					return
				}
				fmt.Fprintf(os.Stderr, "Receive error: %v\n", err)
				return
			}
		}
	}()

	<-ctx.Done()
}
