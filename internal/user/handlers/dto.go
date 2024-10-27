package handlers

import (
	"auth/internal/user/models"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

func UserToModel(u User) models.User {
	return models.User{
		ID:       u.ID,
		Username: u.Email,
	}
}

func UserToDTO(u models.User) User {
	return User{
		ID:    u.ID,
		Email: u.Username,
	}
}
