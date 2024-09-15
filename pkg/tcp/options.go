package tcp

import "time"

const defaultBufferSize = 4 << 10

type ClientOption func(*Client)

func WithClientIdleTimeout(timeout time.Duration) ClientOption {
	return func(client *Client) {
		client.idleTimeout = timeout
	}
}

func WithClientBufferSize(size uint) ClientOption {
	return func(client *Client) {
		client.bufferSize = int(size)
	}
}

type ServerOption func(*Server)

func WithServerIdleTimeout(timeout time.Duration) ServerOption {
	return func(server *Server) {
		server.idleTimeout = timeout
	}
}

func WithServerBufferSize(size uint) ServerOption {
	return func(server *Server) {
		server.bufferSize = int(size)
	}
}

func WithServerMaxConnectionsNumber(count uint) ServerOption {
	return func(server *Server) {
		server.maxConnections = int(count)
	}
}
