package vobjects

import (
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hashPassword []byte
}

func NewPassword(password string) (Password, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}

	return Password{
		hashPassword: hash,
	}, nil
}

func NewExistingPassword(hashPassword string) Password {
	return Password{
		hashPassword: []byte(hashPassword),
	}
}

func (m Password) String() string {
	return "password"
}

func (m Password) Hash() string {
	return string(m.hashPassword)
}

func (m Password) Compare(otherPassword string) bool {
	return bcrypt.CompareHashAndPassword(m.hashPassword, []byte(otherPassword)) == nil
}
