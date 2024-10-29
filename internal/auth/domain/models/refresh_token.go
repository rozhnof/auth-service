package models

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

const tokenLength = 255

type RefreshToken struct {
	Token     string
	ExpiredAt time.Time
	IsRevoked bool
}

func NewRefreshToken(ttl time.Duration) (RefreshToken, error) {
	token, err := generateRandomString(tokenLength)
	if err != nil {
		return RefreshToken{}, err
	}

	rt := RefreshToken{
		Token:     token,
		ExpiredAt: time.Now().Add(ttl),
	}

	return rt, nil
}

func (t *RefreshToken) Valid() bool {
	return t.ExpiredAt.After(time.Now()) && !t.IsRevoked
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
