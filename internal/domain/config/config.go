package config

import "time"

type Config struct {
	Logger *Logger
	Tcp    *TCP
}

type TCP struct {
	Address        string        `yaml:"address"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize string        `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

type Logger struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}
