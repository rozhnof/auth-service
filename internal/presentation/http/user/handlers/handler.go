package http_user_handlers

import (
	user_services "auth/internal/application/services/user"
)

type AuthHandler struct {
	service *user_services.UserService
}

func NewAuthHandler(service *user_services.UserService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}
