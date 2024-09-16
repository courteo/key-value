package initializer

import (
	"context"
)

func (i *Initializer) StartDatabase(ctx context.Context) error {
	i.server.HandleQueries(ctx, func(ctx context.Context, query []byte) []byte {
		response, err := i.db.HandleQuery(ctx, string(query))
		if err != nil {
			return []byte(err.Error())
		}
	
		return []byte(response)
	})

	return nil
}
