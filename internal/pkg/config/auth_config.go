package config

import (
	"auth/pkg/validator"
	"time"
)

type server struct {
	Address         string        `yaml:"address"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type mail struct {
	Host     string `yaml:"host"`
	Sender   string `yaml:"sender"`
	Password string `yaml:"password"`
}

type service struct {
	AccessTokenTimeout  time.Duration `yaml:"access_token_timeout"`
	RefreshTokenTimeout time.Duration `yaml:"refresh_token_timeout"`
}

type AuthConfig struct {
	Server  server  `yaml:"server"`
	Mail    mail    `yaml:"mail"`
	Service service `yaml:"service"`
}

func (c *AuthConfig) Validate() error {
	return validator.ValidYamlFields(c)
}
