package in_memory

import (
	"context"
	"hash/fnv"

	"github.com/courteo/key-value/pkg/common"
	"go.uber.org/zap"
)

func (e *Engine) Get(ctx context.Context, key string) (string, bool) {
	partitionIdx := 0
	if len(e.partitions) > 1 {
		partitionIdx = e.partitionIdx(key)
	}

	partition := e.partitions[partitionIdx]
	value, found := partition.Get(key)

	txID := common.GetTxIDFromContext(ctx)
	e.logger.Debug("successfull get query", zap.Int64("tx", txID))
	return value, found
}

func (e *Engine) Set(ctx context.Context, key string, value string) {
	partitionIdx := 0
	if len(e.partitions) > 1 {
		partitionIdx = e.partitionIdx(key)
	}

	partition := e.partitions[partitionIdx]
	partition.Set(key, value)

	txID := common.GetTxIDFromContext(ctx)
	e.logger.Debug("successfull set query", zap.Int64("tx", txID))
}

func (e *Engine) Delete(ctx context.Context, key string) {
	partitionIdx := 0
	if len(e.partitions) > 1 {
		partitionIdx = e.partitionIdx(key)
	}

	partition := e.partitions[partitionIdx]
	partition.Delete(key)

	txID := common.GetTxIDFromContext(ctx)
	e.logger.Debug("successfull del query", zap.Int64("tx", txID))
}

func (e *Engine) partitionIdx(key string) int {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(key))

	return int(hash.Sum32()) % len(e.partitions)
}
