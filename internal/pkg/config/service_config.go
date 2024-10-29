package config

import "time"

type AccessTokenConfig struct {
	Timeout   time.Duration `yaml:"timeout"   env-required:"true"`
	SecretKey string        `env:"SECRET_KEY" env-required:"true"`
}

type RefreshTokenConfig struct {
	Timeout time.Duration `yaml:"timeout"   env-required:"true"`
}

type TokensConfig struct {
	Access  AccessTokenConfig  `yaml:"access"  env-required:"true"`
	Refresh RefreshTokenConfig `yaml:"refresh" env-required:"true"`
}

type ServiceConfig struct {
	Tokens TokensConfig `yaml:"tokens"`
}
