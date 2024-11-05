package http_handlers

import (
	"auth/internal/auth/application/services"
	"log/slog"
	"net/http"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Login @Summary User login
// @Description Login user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login Request"
// @Success 200 {object} LoginResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	ctx, span := otel.Tracer("auth-handler").Start(c.Request.Context(), "AuthHandler.Login")
	defer span.End()

	log := h.log.With(
		slog.String("function", "AuthHandler.Login"),
	)

	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Info("bad request")

		c.String(http.StatusBadRequest, err.Error())
		return
	}

	log = log.With(
		slog.String("username", request.Username),
	)

	at, rt, err := h.userService.Login(ctx, request.Username, request.Password)
	if err != nil {
		log.Info("failed user login", slog.String("error", err.Error()))

		if errors.Is(err, services.ErrInvalidPassword) {
			c.String(http.StatusOK, "invalid username or password")
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	response := LoginResponse{
		AccessToken:  at,
		RefreshToken: rt,
	}

	c.JSON(http.StatusOK, response)
}
