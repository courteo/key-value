package config

type Config struct{}

type Logger struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}
