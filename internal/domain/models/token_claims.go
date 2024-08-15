package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenClaims struct {
	TokenID uuid.UUID
	UserID  uuid.UUID
	UserIP  string
	jwt.RegisteredClaims
}
