package domain

import "golang.org/x/exp/rand"

type Secret []byte

func (s Secret) String() string {
	return "secret"
}

func (s Secret) Get() []byte {
	return s
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = letters[rand.Intn(len(letters))]
	}

	return string(bytes)
}
