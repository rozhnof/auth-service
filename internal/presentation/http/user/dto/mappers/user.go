package http_user_mappers

import (
	"auth/internal/domain/models"
	http_user_structs "auth/internal/presentation/http/user/dto/structs"
)

func UserToModel(u http_user_structs.User) models.User {
	return models.User{
		ID:       u.ID,
		Username: u.Email,
	}
}

func UserToDTO(u models.User) http_user_structs.User {
	return http_user_structs.User{
		ID:    u.ID,
		Email: u.Username,
	}
}
