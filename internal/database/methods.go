package database

import (
	"context"
	"fmt"

	"github.com/courteo/key-value/internal/domain"
	"github.com/courteo/key-value/internal/domain/command"
)

func (d *Database) HandleQuery(ctx context.Context, queryStr string) (string, error) {
	query, err := d.Computer.ParseQuery(queryStr)
	if err != nil {
		return "", err
	}

	var value string

	switch query.Command {
	case command.SetID:
		value, err = d.handleSetQuery(ctx, query)
	case command.GetID:
		value, err = d.handleGetQuery(ctx, query)
	case command.DeleteID:
		value, err = d.handleDeleteQuery(ctx, query)
	}

	if err != nil {
		return "", err
	}

	return value, nil
}

func (d *Database) handleSetQuery(ctx context.Context, query domain.Query) (string, error) {
	err := d.Storage.Set(ctx, query.Key, query.Value)
	if err != nil {
		return "", err
	}
	return "[ok]", nil
}

func (d *Database) handleDeleteQuery(ctx context.Context, query domain.Query) (string, error) {
	err := d.Storage.Delete(ctx, query.Key)
	if err != nil {
		return "", err
	}
	return "[ok]", nil
}

func (d *Database) handleGetQuery(ctx context.Context, query domain.Query) (string, error) {
	val, err := d.Storage.Get(ctx, query.Key)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("[ok] %s", val), nil
}
