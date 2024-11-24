package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/courteo/key-value/internal/domain/config"
	initializer2 "github.com/courteo/key-value/internal/initializer"
	config2 "github.com/courteo/key-value/pkg/config"
)

var (
	ConfigFileName = os.Getenv("CONFIG_FILE_NAME")
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := &config.Config{}
	data, err := os.ReadFile("../../config.yml")
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(data)
	cfg, err = config2.New(reader)
	if err != nil {
		log.Fatal(err)
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
