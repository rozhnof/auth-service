package models

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type AccessToken struct {
	Token string
}

func NewAccessToken(ttl time.Duration, secretKey []byte) (AccessToken, error) {
	expiredAt := time.Now().Add(ttl)

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiredAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "auth-service",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return AccessToken{}, err
	}

	at := AccessToken{
		Token: signedToken,
	}

	return at, nil
}
