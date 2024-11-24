package tcp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type Client struct {
	connection  net.Conn
	idleTimeout time.Duration
	bufferSize  int
}

func NewClient(address string, options ...ClientOption) (*Client, error) {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	client := &Client{
		connection: connection,
		bufferSize: defaultBufferSize,
	}

	for _, option := range options {
		option(client)
	}

	if client.idleTimeout != 0 {
		err = connection.SetDeadline(time.Now().Add(client.idleTimeout))
		if err != nil {
			return nil, fmt.Errorf("failed to set deadline for connection: %w", err)
		}
	}

	return client, nil
}

func (c *Client) Send(request []byte) ([]byte, error) {
	_, err := c.connection.Write(request)
	if err != nil {
		return nil, err
	}

	response := make([]byte, c.bufferSize)

	count, err := c.connection.Read(response)
	if err != nil && err != io.EOF {
		return nil, err
	}

	if count >= c.bufferSize {
		return nil, errors.New("small buffer size")
	}

	return response[:count], nil
}

func (c *Client) Close() {
	if c.connection != nil {
		_ = c.connection.Close()
	}
}
