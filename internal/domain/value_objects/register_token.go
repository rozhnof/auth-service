package vobjects

import (
	"time"

	"github.com/rozhnof/auth-service/internal/domain"
)

const (
	registerTokenLength = 25
	registerTokenTTL    = time.Hour * 24
)

type RegisterToken struct {
	token     string
	expiredAt time.Time
}

func NewRegisterToken() *RegisterToken {
	return &RegisterToken{
		token:     domain.GenerateRandomString(registerTokenLength),
		expiredAt: time.Now().Add(registerTokenTTL),
	}
}

func NewExistingRegisterToken(token string, expiredAt time.Time) *RegisterToken {
	return &RegisterToken{
		token:     token,
		expiredAt: expiredAt,
	}
}

func (t RegisterToken) Token() string {
	return t.token
}

func (t RegisterToken) ExpiredAt() time.Time {
	return t.expiredAt
}

func (t RegisterToken) Valid() bool {
	return t.expiredAt.After(time.Now())
}
