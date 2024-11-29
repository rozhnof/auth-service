package domain

import (
	"crypto/rand"
	"encoding/base64"
)

type Secret []byte

func (s Secret) String() string {
	return "secret"
}

func (s Secret) Get() []byte {
	return s
}

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
