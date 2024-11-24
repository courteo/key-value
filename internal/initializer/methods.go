package initializer

import (
	"context"

	"github.com/courteo/key-value/internal/database"
	"github.com/courteo/key-value/internal/database/compute"
	"github.com/courteo/key-value/internal/database/storage"
	"golang.org/x/sync/errgroup"
)

func (i *Initializer) StartDatabase(ctx context.Context) error {
	compute := compute.New(i.logger)

	var options []storage.StorageOption
	if i.wal != nil {
		options = append(options, storage.WithWAL(i.wal))
	}

	if i.master != nil {
		options = append(options, storage.WithReplication(i.master))
	} else if i.slave != nil {
		options = append(options, storage.WithReplication(i.slave))
		options = append(options, storage.WithReplicationStream(i.slave.ReplicationStream()))
	}

	storage := storage.New(i.logger, i.engine, options...)

	database := database.New(i.logger, compute, storage)

	group, groupCtx := errgroup.WithContext(ctx)
	if i.wal != nil {
		if i.slave != nil {
			group.Go(func() error {
				i.slave.Start(groupCtx)
				return nil
			})
		} else {
			group.Go(func() error {
				i.wal.Start(groupCtx)
				return nil
			})
		}

		if i.master != nil {
			group.Go(func() error {
				i.master.Start(groupCtx)
				return nil
			})
		}
	}

	group.Go(func() error {
		i.server.HandleQueries(groupCtx, func(ctx context.Context, query []byte) []byte {
			response, err := database.HandleQuery(ctx, string(query))
			if err != nil {
				return []byte(err.Error())
			}

			return []byte(response)
		})

		return nil
	})

	return group.Wait()
}
