package in_memory

import (
	"context"
)

func (e *Engine) Get(ctx context.Context, key string) (string, bool) {
	value, ok := e.hashTable.Get(key)

	//e.logger.Debug("get key successfully", zap.Int64("request_id", ctx.Value("request_id").(int64)))

	return value, ok
}

func (e *Engine) Set(ctx context.Context, key string, value string) {
	e.hashTable.Set(key, value)

	//e.logger.Debug("set key successfully", zap.Int64("request_id", ctx.Value("request_id").(int64)))
}

func (e *Engine) Delete(ctx context.Context, key string) {
	e.hashTable.Delete(key)

	//e.logger.Debug("del key successfully", zap.Int64("request_id", ctx.Value("request_id").(int64)))
}
