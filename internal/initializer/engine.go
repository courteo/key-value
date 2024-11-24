package initializer

import (
	"errors"

	in_memory "github.com/courteo/key-value/internal/database/storage/engine/in-memory"
	"github.com/courteo/key-value/internal/domain/config"
	"go.uber.org/zap"
)

func createEngine(cfg *config.Engine, logger *zap.Logger) (*in_memory.Engine, error) {
	if cfg.Type != "" {
		supportedTypes := map[string]struct{}{
			"in_memory": {},
		}

		if _, found := supportedTypes[cfg.Type]; !found {
			return nil, errors.New("engine type is incorrect")
		}
	}

	var options []in_memory.EngineOption
	if cfg.PartitionsNumber != 0 {
		options = append(options, in_memory.WithPartitions(cfg.PartitionsNumber))
	}

	engine := in_memory.New(logger, options...)

	return engine, nil
}
