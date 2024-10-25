package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server     ServerConfig `yaml:"server" env-required:"true"`
	Repository RepositoryConfig
	Logger     LoggerConfig `yaml:"logging" env-required:"true"`
}

func NewConfig(configPath string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
