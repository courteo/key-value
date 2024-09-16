package database

import (
	"context"
	"github.com/courteo/key-value/internal/domain"
	"go.uber.org/zap"
)

type Storage interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) string
	Delete(ctx context.Context, key string) string
}

type Computer interface {
	ParseQuery(query string) (domain.Query, error)
}

type Database struct {
	Computer Computer
	Storage  Storage
	logger   *zap.Logger
}

func New(logger *zap.Logger, computer Computer, storage Storage) *Database {
	return &Database{
		Computer: computer,
		Storage:  storage,
		logger:   logger,
	}
}
