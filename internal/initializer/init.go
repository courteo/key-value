package initializer

import (
	"errors"
	"github.com/courteo/key-value/internal/database"
	"github.com/courteo/key-value/internal/database/compute"
	"github.com/courteo/key-value/internal/database/storage"
	in_memory "github.com/courteo/key-value/internal/database/storage/engine/in-memory"
	"github.com/courteo/key-value/internal/domain/config"
	logger2 "github.com/courteo/key-value/pkg/logger"
	"github.com/courteo/key-value/pkg/tcp"
	"go.uber.org/zap"
)

const defaultAddress = ":3000"

type Initializer struct {
	db     *database.Database
	server *tcp.Server
	logger *zap.Logger
}

func New(cfg *config.Config) (*Initializer, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	logger, err := logger2.New(cfg.Logger)
	if err != nil {
		return nil, err
	}

	address := getAddress(cfg.Tcp)

	server, err := tcp.NewServer(address, logger)
	if err != nil {
		return nil, err
	}

	computer := compute.New(logger)

	engine := in_memory.New(logger)
	storage := storage.New(logger, engine)

	db := database.New(logger, computer, storage)

	return &Initializer{
		db:     db,
		logger: logger,
		server: server,
	}, nil
}

func getAddress(cfg *config.TCP) string {
	if cfg == nil || cfg.Address == "" {
		return defaultAddress
	}

	return cfg.Address
}
