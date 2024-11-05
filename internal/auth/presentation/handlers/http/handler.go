package http_handlers

import (
	"auth/internal/auth/application/services"
	"log/slog"
)

type AuthHandler struct {
	log         *slog.Logger
	userService *services.UserService
}

func NewAuthHandler(service *services.UserService, log *slog.Logger) *AuthHandler {
	log = log.With(
		slog.String("layer", "presentation"),
		slog.String("pkg", "http_handlers"),
	)

	return &AuthHandler{
		userService: service,
		log:         log,
	}
}

const tracerName = "Auth Service"
