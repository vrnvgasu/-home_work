package main

import (
	"fmt"
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

type client struct {
	io.Reader
	io.Writer
	conn    net.Conn
	address string
	timeout time.Duration
}

func (c *client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("client Connect: DialTimeout failed: %w", err)
	}
	c.conn = conn

	return nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) Send() error {
	if _, err := io.Copy(c.conn, c.Reader); err != nil {
		return fmt.Errorf("client Send: Copy failed: %w", err)
	}

	return nil
}

func (c *client) Receive() error {
	if _, err := io.Copy(c.Writer, c.conn); err != nil {
		return fmt.Errorf("client Receive: Copy failed: %w", err)
	}

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		Reader:  in,
		Writer:  out,
		address: address,
		timeout: timeout,
	}
}
