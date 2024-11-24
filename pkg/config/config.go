package config

import (
	"errors"
	"fmt"
	"io"

	"github.com/courteo/key-value/internal/domain/config"
	"gopkg.in/yaml.v3"
)

func New(reader io.Reader) (*config.Config, error) {
	if reader == nil {
		return nil, errors.New("incorrect reader")
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.New("falied to read buffer")
	}

	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
