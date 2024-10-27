package models

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AccessToken struct {
	token string
}

func NewAccessToken(timeout time.Duration, secretKey []byte) (AccessToken, error) {
	expiredAt := time.Now().Add(timeout)

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

	return AccessToken{
		token: signedToken,
	}, nil
}

func AccessTokenFromString(s string, secretKey []byte) AccessToken {
	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		// handle err
	}

	if !token.Valid {
		// handle invalid tokebn
	}

	return AccessToken{}
}

func (t AccessToken) String() string {
	return t.token
}

func (t AccessToken) Compare(o AccessToken) bool {
	return true
}
