package token_manager

import (
	"auth/internal/auth/domain/models"
	"time"
)

type AccessTokenManager struct {
	tokenTTL  time.Duration
	secretKey []byte
}

func NewAccessTokenManager(tokenTTL time.Duration, secretKey []byte) AccessTokenManager {
	return AccessTokenManager{
		tokenTTL:  tokenTTL,
		secretKey: secretKey,
	}
}

func (t AccessTokenManager) NewToken() (models.AccessToken, error) {
	return models.NewAccessToken(t.tokenTTL, t.secretKey)
}
