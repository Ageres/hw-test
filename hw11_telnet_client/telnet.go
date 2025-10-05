package main

import (
	"io"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
}

// Close implements TelnetClient.
func (t *telnetClient) Close() error {
	panic("unimplemented")
}

// Connect implements TelnetClient.
func (t *telnetClient) Connect() error {
	panic("unimplemented")
}

// Receive implements TelnetClient.
func (t *telnetClient) Receive() error {
	panic("unimplemented")
}

// Send implements TelnetClient.
func (t *telnetClient) Send() error {
	panic("unimplemented")
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	// Place your code here.

	client := telnetClient{
		address: address,
		timeout: timeout,
	}

	return &client
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
