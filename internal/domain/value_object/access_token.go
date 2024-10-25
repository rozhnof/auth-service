package value_object

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AccessToken struct {
	token *jwt.Token
}

func NewAccessToken(timeout time.Duration) AccessToken {
	expiredAt := time.Now().Add(timeout)

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiredAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "auth-service",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return AccessToken{
		token: token,
	}
}

func (t AccessToken) Sign(secretKey []byte) (string, error) {
	return t.token.SignedString(secretKey)
}
