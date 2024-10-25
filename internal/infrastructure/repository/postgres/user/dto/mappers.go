package postgres_user_dto

import (
	"auth/internal/domain/models"
	"auth/internal/domain/value_object"
)

func UserToModel(u *User) *models.User {
	return &models.User{
		ID:           u.ID,
		Username:     u.Username,
		Password:     value_object.PasswordFromString(u.Password),
		RefreshToken: value_object.RefreshTokenFromString(u.RefreshToken),
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
