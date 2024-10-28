package token_manager

import (
	"crypto/rand"
	"encoding/base64"
)

const tokenLength = 77

type RefreshTokenManager struct{}

func NewRefreshTokenManager(secretKey []byte) *RefreshTokenManager {
	return &RefreshTokenManager{}
}

func (t RefreshTokenManager) NewRefreshToken() (string, error) {
	return generateRandomString(tokenLength)
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
