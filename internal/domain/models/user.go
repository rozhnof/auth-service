package models

import (
	"auth/internal/domain/value_object"
)

type User struct {
	ID           value_object.UserID
	Username     value_object.Username
	Password     value_object.Password
	AccessToken  value_object.AccessToken
	RefreshToken value_object.RefreshToken
}
