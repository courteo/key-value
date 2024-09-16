package main

import (
	"bytes"
	"context"
	"github.com/courteo/key-value/internal/domain/config"
	initializer2 "github.com/courteo/key-value/internal/initializer"
	config2 "github.com/courteo/key-value/pkg/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	ConfigFileName = os.Getenv("CONFIG_FILE_NAME")
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := &config.Config{}
	if ConfigFileName != "" {
		data, err := os.ReadFile(ConfigFileName)
		if err != nil {
			log.Fatal(err)
		}

		reader := bytes.NewReader(data)
		cfg, err = config2.New(reader)
		if err != nil {
			log.Fatal(err)
		}
	}

	initializer, err := initializer2.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = initializer.StartDatabase(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
