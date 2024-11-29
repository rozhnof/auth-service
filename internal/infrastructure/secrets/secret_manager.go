package secrets

import (
	"os"

	"github.com/rozhnof/auth-service/internal/domain"
)

const (
	envSecretKey          = "SECRET_KEY"
	envGoogleClientID     = "GOOGLE_CLIENT_ID"
	envGoogleClientSecret = "GOOGLE_CLIENT_SECRET"
)

type EnvSecretManager struct{}

func NewEnvSecretManager() EnvSecretManager {
	return EnvSecretManager{}
}

func (m EnvSecretManager) SecretKey() domain.Secret {
	return domain.Secret(os.Getenv(envSecretKey))
}

func (m EnvSecretManager) GoogleClientID() domain.Secret {
	return domain.Secret(os.Getenv(envGoogleClientID))
}

func (m EnvSecretManager) GoogleClientSecret() domain.Secret {
	return domain.Secret(os.Getenv(envGoogleClientSecret))
}
