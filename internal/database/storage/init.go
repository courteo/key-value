package storage

import (
	"context"
	"go.uber.org/zap"
)

type Engine interface {
	Get(ctx context.Context, key string) (string, bool)
	Set(ctx context.Context, key string, value string)
	Delete(ctx context.Context, key string)
}

type Storage struct {
	Engine Engine
	logger *zap.Logger
}

func New(logger *zap.Logger, engine Engine) *Storage {
	return &Storage{
		Engine: engine,
		logger: logger,
	}
}
