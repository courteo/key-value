package in_memory

import (
	hash_table "github.com/courteo/key-value/internal/domain/hash-table"
	"go.uber.org/zap"
)

type Engine struct {
	hashTable *hash_table.HashTable
	logger    *zap.Logger
}

func New(logger *zap.Logger) *Engine {
	return &Engine{
		hashTable: hash_table.New(),
		logger:    logger,
	}
}
