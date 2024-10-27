package handlers

import "auth/internal/user/services"

type AuthHandler struct {
	service *services.UserService
}

func NewAuthHandler(service *services.UserService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}
