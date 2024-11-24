package in_memory

import hash_table "github.com/courteo/key-value/internal/domain/hash-table"

type EngineOption func(*Engine)

func WithPartitions(partitionsNumber uint) EngineOption {
	return func(engine *Engine) {
		engine.partitions = make([]*hash_table.HashTable, partitionsNumber)
		for i := 0; i < int(partitionsNumber); i++ {
			engine.partitions[i] = hash_table.New()
		}
	}
}
