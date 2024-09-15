package compute

import (
	"go.uber.org/zap"
)

type Computer struct {
	logger *zap.Logger
}

func NewComputer(logger *zap.Logger) *Computer {
	return &Computer{
		logger: logger,
	}
}
