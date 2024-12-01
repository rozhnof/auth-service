package vobjects

import (
	"time"

	"github.com/rozhnof/auth-service/internal/domain"
)

const (
	refreshTokenLength = 255
)

type RefreshToken struct {
	token     string
	expiredAt time.Time
}

func NewRefreshToken(ttl time.Duration) *RefreshToken {
	return &RefreshToken{
		token:     domain.GenerateRandomString(refreshTokenLength),
		expiredAt: time.Now().Add(ttl),
	}
}

func NewExistingRefreshToken(token string, expiredAt time.Time) *RefreshToken {
	return &RefreshToken{
		token:     token,
		expiredAt: expiredAt,
	}
}

func (t RefreshToken) Token() string {
	return t.token
}

func (t RefreshToken) ExpiredAt() time.Time {
	return t.expiredAt
}

func (t RefreshToken) Compare(other RefreshToken) bool {
	return t.token == other.token
}

func (t RefreshToken) Valid() bool {
	return t.expiredAt.After(time.Now())
}
