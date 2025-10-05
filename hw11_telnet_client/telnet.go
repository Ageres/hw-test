package main

import (
	"io"
	"net"
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
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	client := telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}

	return &client
}

func (t *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}
	t.conn = conn
	return nil
}

// Close implements TelnetClient.
func (t *telnetClient) Close() error {
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

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
