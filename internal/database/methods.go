package database

import (
	"context"
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
	case command.Set:
		err = d.handleSetQuery(ctx, query)
	case command.Get:
		value, err = d.handleGetQuery(ctx, query)
	case command.Delete:
		err = d.handleDeleteQuery(ctx, query)
	}

	if err != nil {
		return "", err
	}

	return value, nil
}

func (d *Database) handleSetQuery(ctx context.Context, query domain.Query) error {
	err := d.Storage.Set(ctx, query.Key, query.Value)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) handleDeleteQuery(ctx context.Context, query domain.Query) error {
	err := d.Storage.Delete(ctx, query.Key)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) handleGetQuery(ctx context.Context, query domain.Query) (string, error) {
	val, err := d.Storage.Get(ctx, query.Key)
	if err != nil {
		return "", err
	}

	return val, nil
}
