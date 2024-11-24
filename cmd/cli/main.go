package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/courteo/key-value/pkg/tcp"
	"go.uber.org/zap"
	"os"
	"syscall"
	"time"
)

func main() {
	address := flag.String("address", "localhost:3223", "Address of the spider")
	idleTimeout := flag.Duration("idle_timeout", 50*time.Minute, "Idle timeout for connection")
	//maxMessageSizeStr := flag.String("max_message_size", "4KB", "Max message size for connection")
	flag.Parse()

	logger, _ := zap.NewProduction()
	//maxMessageSize, err := common.ParseSize(*maxMessageSizeStr)
	//if err != nil {
	//	logger.Fatal("failed to parse max message size", zap.Error(err))
	//}

	var options []tcp.ClientOption
	options = append(options, tcp.WithClientIdleTimeout(*idleTimeout))
	//options = append(options, tcp.WithClientBufferSize(uint(maxMessageSize)))

	reader := bufio.NewReader(os.Stdin)
	client, err := tcp.NewClient(*address, options...)
	if err != nil {
		logger.Fatal("failed to connect with server", zap.Error(err))
	}

	for {
		fmt.Print("[db] > ")
		
		request, err := reader.ReadString('\n')
		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal("connection was closed", zap.Error(err))
		} else if err != nil {
			logger.Error("failed to read query", zap.Error(err))
		}

		response, err := client.Send([]byte(request))
		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal("connection was closed", zap.Error(err))
		} else if err != nil {
			logger.Error("failed to send query", zap.Error(err))
		}

		fmt.Println(string(response))
	}
}
