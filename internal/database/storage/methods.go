package storage

import (
	"context"
	"errors"
)

var errNotFound = errors.New("key not found")

func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	value, ok := s.Engine.Get(ctx, key)
	if !ok {
		return "", errNotFound
	}

	return value, nil
}

func (s *Storage) Set(ctx context.Context, key string, value string) string {
	s.Engine.Set(ctx, key, value)

	return "was set"
}

func (s *Storage) Delete(ctx context.Context, key string) string {
	s.Engine.Delete(ctx, key)

	return "was delete"
}
