package handlers

import "auth/internal/domain/services"

const queryParamUserID = "user_id"

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}
