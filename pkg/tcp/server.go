package tcp

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"sync"
	"time"
)

type Handler = func(context.Context, []byte) []byte

type Server struct {
	listener net.Listener

	idleTimeout    time.Duration
	bufferSize     int
	maxConnections int

	logger *zap.Logger
}

func NewServer(address string, logger *zap.Logger, options ...ServerOption) (*Server, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	server := &Server{
		listener: listener,
		logger:   logger,
	}

	for _, option := range options {
		option(server)
	}

	if server.bufferSize == 0 {
		server.bufferSize = defaultBufferSize
	}

	return server, nil
}

func (s *Server) HandleQueries(ctx context.Context, handler Handler) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			connection, err := s.listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				s.logger.Error("failed to accept", zap.Error(err))
				continue
			}

			wg.Add(1)

			go func(connection net.Conn) {
				defer func() {
					wg.Done()
				}()

				s.handleConnection(ctx, connection, handler)
			}(connection)
		}
	}()

	go func() {
		defer wg.Done()

		<-ctx.Done()
		s.listener.Close()
	}()

	wg.Wait()
}

func (s *Server) handleConnection(ctx context.Context, connection net.Conn, handler Handler) {
	defer func() {
		if v := recover(); v != nil {
			s.logger.Error("captured panic", zap.Any("panic", v))
		}

		if err := connection.Close(); err != nil {
			s.logger.Warn("failed to close connection", zap.Error(err))
		}
	}()

	request := make([]byte, s.bufferSize)

	for {
		if s.idleTimeout != 0 {
			err := connection.SetDeadline(time.Now().Add(s.idleTimeout))
			if err != nil {
				s.logger.Warn("failed to set read deadline", zap.Error(err))
				break
			}
		}

		count, err := connection.Read(request)
		if err != nil && err != io.EOF {
			s.logger.Warn("failed to read", zap.Error(err))
			break
		}

		if count == s.bufferSize {
			s.logger.Warn("small buffer size")
			break
		}

		response := handler(ctx, request[:count])

		_, err = connection.Write(response)
		if err != nil {
			s.logger.Warn("failed to write", zap.Error(err))
			break
		}
	}
}
