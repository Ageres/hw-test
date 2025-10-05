package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var (
	timeout time.Duration
)

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?

	ctx := context.Background()

	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		log.Fatalf("Usage: go-telnet [--timeout=10s] host port")
	}

	host := args[0]
	port := args[1]

	log.Printf("timeout '%s', host '%s', port '%s", timeout, host, port)

	address := net.JoinHostPort(host, port)

	log.Printf("address '%s'", address)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatalf("connect error: %v", err)
	}
	defer client.Close()

	log.Printf("connected to %s\n", address)

	go func() {
		for {
			if err := client.Receive(); err != nil {
				if err != io.EOF {
					fmt.Fprintf(os.Stderr, "receive error: %v\n", err)
				}
				return
			}
		}
	}()

	go func() {
		for {
			if err := client.Send(); err != nil {
				if err != io.EOF {
					fmt.Fprintf(os.Stderr, "send error: %v\n", err)
				}
				return
			}
		}
	}()

	<-ctx.Done()

}
