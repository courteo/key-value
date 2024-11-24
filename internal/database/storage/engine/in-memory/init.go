package in_memory

import (
	hash_table "github.com/courteo/key-value/internal/domain/hash-table"
	"go.uber.org/zap"
)

type Engine struct {
	partitions []*hash_table.HashTable
	logger     *zap.Logger
}

func New(logger *zap.Logger, options ...EngineOption) *Engine {
	engine := &Engine{
		logger: logger,
	}

	for _, option := range options {
		option(engine)
	}

	if len(engine.partitions) == 0 {
		engine.partitions = make([]*hash_table.HashTable, 1)
		engine.partitions[0] = hash_table.New()
	}

	return engine
}
