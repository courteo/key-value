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

func (s *Storage) Set(ctx context.Context, key string, value string) error {
	s.Engine.Set(ctx, key, value)

	return nil
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	s.Engine.Delete(ctx, key)

	return nil
}
