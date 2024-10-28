package handlers

import "auth/internal/auth/application/services"

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(service *services.UserService) *AuthHandler {
	return &AuthHandler{
		userService: service,
	}
}
