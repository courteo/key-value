package storage

import (
	"context"
	"errors"

	"github.com/courteo/key-value/internal/database/storage/wal"
	"github.com/courteo/key-value/internal/domain/command"
	"github.com/courteo/key-value/pkg/common"
)

var (
	ErrorNotFound  = errors.New("not found")
	ErrorMutableTX = errors.New("mutable transaction on slave")
)

func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	txID := s.generator.Generate()
	ctx = common.ContextWithTxID(ctx, txID)

	value, found := s.engine.Get(ctx, key)
	if !found {
		return "", ErrorNotFound
	}

	return value, nil
}

func (s *Storage) Set(ctx context.Context, key string, value string) error {
	if s.replica != nil && !s.replica.IsMaster() {
		return ErrorMutableTX
	} else if ctx.Err() != nil {
		return ctx.Err()
	}

	txID := s.generator.Generate()
	ctx = common.ContextWithTxID(ctx, txID)

	if s.wal != nil {
		futureResponse := s.wal.Set(ctx, key, value)
		if err := futureResponse.Get(); err != nil {
			return err
		}
	}

	s.engine.Set(ctx, key, value)
	return nil
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	if s.replica != nil && !s.replica.IsMaster() {
		return ErrorMutableTX
	} else if ctx.Err() != nil {
		return ctx.Err()
	}

	txID := s.generator.Generate()
	ctx = common.ContextWithTxID(ctx, txID)

	if s.wal != nil {
		futureResponse := s.wal.Del(ctx, key)
		if err := futureResponse.Get(); err != nil {
			return err
		}
	}

	s.engine.Delete(ctx, key)
	return nil
}

func (s *Storage) applyData(logs []wal.Log) int64 {
	var lastLSN int64
	for _, log := range logs {
		lastLSN = max(lastLSN, log.LSN)
		ctx := common.ContextWithTxID(context.Background(), log.LSN)
		switch log.CommandID {
		case command.SetID:
			s.engine.Set(ctx, log.Arguments[0], log.Arguments[1])
		case command.DeleteID:
			s.engine.Delete(ctx, log.Arguments[0])
		}
	}

	return lastLSN
}
