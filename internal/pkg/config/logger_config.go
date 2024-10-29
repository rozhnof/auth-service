package config

type LoggerConfig struct {
	Level string `yaml:"level"`
	Path  string `yaml:"path"`
}
