package config

import "time"

type Tokens struct {
	AccessTokenTTL  time.Duration `yaml:"access_ttl"  env-required:"true"`
	RefreshTokenTTL time.Duration `yaml:"refresh_ttl" env-required:"true"`
}
