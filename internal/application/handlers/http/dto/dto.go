package dto

import (
	"auth/internal/domain/models"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

func (u User) ToModel() models.User {
	return models.User{
		ID:    u.ID,
		Email: u.Email,
	}
}

func UserToDTO(u models.User) User {
	return User{
		ID:    u.ID,
		Email: u.Email,
	}
}
