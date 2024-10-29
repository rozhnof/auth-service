package token_manager

import (
	"auth/internal/auth/domain/models"
	"time"
)

type RefreshTokenManager struct {
	tokenTTL time.Duration
}

func NewRefreshTokenManager(tokenTTL time.Duration) *RefreshTokenManager {
	return &RefreshTokenManager{
		tokenTTL: tokenTTL,
	}
}

func (t RefreshTokenManager) NewToken() (models.RefreshToken, error) {
	return models.NewRefreshToken(t.tokenTTL)
}
