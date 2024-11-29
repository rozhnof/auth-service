package vobjects

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const Issuer = "stakewolle-auth-service"

type AccessTokenPayload struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	AccessTokenPayload
}

type AccessToken struct {
	token string
}

func NewAccessToken(ttl time.Duration, secretKey []byte, payload AccessTokenPayload) (AccessToken, error) {
	claims := AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: Issuer,
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(ttl),
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
		},
		AccessTokenPayload: payload,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return AccessToken{}, err
	}

	at := AccessToken{
		token: signedToken,
	}

	return at, nil
}

func (t AccessToken) Token() string {
	return t.token
}

func VerifyAccessToken(tokenString string, secretKey []byte) (AccessTokenClaims, error) {
	var claims AccessTokenClaims

	if _, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	}); err != nil {
		return AccessTokenClaims{}, err
	}

	return claims, nil
}
