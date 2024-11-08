package http_handlers

import (
	"auth/internal/auth/application/services"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type AuthHandler struct {
	log         *slog.Logger
	userService *services.UserService
	tracer      trace.Tracer
}

func NewAuthHandler(service *services.UserService, log *slog.Logger, tracer trace.Tracer) *AuthHandler {
	return &AuthHandler{
		userService: service,
		log:         log,
		tracer:      tracer,
	}
}
