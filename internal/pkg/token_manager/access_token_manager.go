package token_manager

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AccessTokenManager struct {
	secretKey []byte
}

func NewAccessTokenManager(secretKey []byte) *AccessTokenManager {
	return &AccessTokenManager{
		secretKey: secretKey,
	}
}

func (t AccessTokenManager) NewAccessTokenWithTTL(ttl time.Duration) (string, error) {
	expiredAt := time.Now().Add(ttl)

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiredAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "auth-service",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString(t.secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
