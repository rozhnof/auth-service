package config

import "time"

type HTTPServerConfig struct {
	Address         string        `yaml:"address"          env-required:"true"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-required:"true"`
}

type ServerConfig struct {
	HTTP HTTPServerConfig `yaml:"http" env-required:"true"`
}
