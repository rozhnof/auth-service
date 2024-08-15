package dao

import (
	"auth/internal/domain/models"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	TokenHash []byte    `db:"token_hash"`
}

func (u User) ToModel() models.User {
	return models.User{
		ID:    u.ID,
		Email: u.Email,
	}
}
