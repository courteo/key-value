package storage

import (
	"context"

	"github.com/courteo/key-value/internal/database/storage/wal"
	"github.com/courteo/key-value/pkg/concurrency"
	"go.uber.org/zap"
)

type Engine interface {
	Get(ctx context.Context, key string) (string, bool)
	Set(ctx context.Context, key string, value string)
	Delete(ctx context.Context, key string)
}

type WAL interface {
	Recover() ([]wal.Log, error)
	Set(context.Context, string, string) concurrency.FutureError
	Del(context.Context, string) concurrency.FutureError
}

type Replica interface {
	IsMaster() bool
}

type Storage struct {
	engine    Engine
	replica   Replica
	wal       WAL
	stream    <-chan []wal.Log
	generator *IDGenerator
	logger    *zap.Logger
}

func New(logger *zap.Logger, engine Engine, options ...StorageOption) *Storage {
	storage := &Storage{
		engine: engine,
		logger: logger,
	}

	for _, option := range options {
		option(storage)
	}

	var lastLSN int64
	if storage.wal != nil {
		logs, err := storage.wal.Recover()
		if err != nil {
			logger.Error("failed to recover data from WAL", zap.Error(err))
		} else {
			lastLSN = storage.applyData(logs)
		}
	}

	if storage.stream != nil {
		go func() {
			for logs := range storage.stream {
				_ = storage.applyData(logs)
			}
		}()
	}

	storage.generator = NewIDGenerator(lastLSN)

	return storage
}
