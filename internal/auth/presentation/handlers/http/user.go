package http_handlers

import (
	"auth/internal/auth/domain/models"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func UserToModel(u User) models.User {
	return models.User{
		ID:       u.ID,
		Username: u.Username,
	}
}

func UserToDTO(u models.User) User {
	return User{
		ID:       u.ID,
		Username: u.Username,
	}
}
