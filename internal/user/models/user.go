package models

import (
	"github.com/google/uuid"
)

type (
	UserID   = uuid.UUID
	Username = string
)

type User struct {
	ID           UserID
	Username     Username
	Password     UserPassword
	AccessToken  AccessToken
	RefreshToken RefreshToken
}
