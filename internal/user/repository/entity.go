package repository

import (
	"auth/internal/user/models"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Username     string    `db:"username"`
	Password     string    `db:"password"`
	RefreshToken string    `db:"refresh_token"`
}

func UserToModel(u *User) *models.User {
	return &models.User{
		ID:           u.ID,
		Username:     u.Username,
		Password:     models.UserPasswordFromString(u.Password),
		RefreshToken: models.RefreshTokenFromString(u.RefreshToken),
	}
}

func UserListToModel(entityUserList []User) []models.User {
	userList := make([]models.User, 0, len(entityUserList))
	for _, userEntity := range entityUserList {
		userList = append(userList, *UserToModel(&userEntity))
	}

	return userList
}

func UserToEntity(u *models.User) *User {
	return &User{
		ID:           u.ID,
		Username:     u.Username,
		Password:     u.Password.String(),
		RefreshToken: u.RefreshToken.String(),
	}
}

func UserListToEntity(entityUserList []models.User) []User {
	userList := make([]User, 0, len(entityUserList))
	for _, userEntity := range entityUserList {
		userList = append(userList, *UserToEntity(&userEntity))
	}

	return userList
}
