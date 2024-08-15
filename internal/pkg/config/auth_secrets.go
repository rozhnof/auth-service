package config

import "auth/pkg/validator"

type AuthSecrets struct {
	DatabaseURL string `env:"DATABASE_URL"`
	SecretKey   string `env:"SECRET_KEY"`
}

func (s *AuthSecrets) Validate() error {
	return validator.ValidEnvFields(s)
}
