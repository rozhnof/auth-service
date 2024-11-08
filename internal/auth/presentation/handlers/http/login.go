//go:generate ffjson $GOFILE
package http_handlers

import (
	"auth/internal/auth/application/services"
	"net/http"

	"github.com/pkg/errors"

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
	ctx, span := h.tracer.Start(c.Request.Context(), "AuthHandler.Login")
	defer span.End()

	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	at, rt, err := h.userService.Login(ctx, request.Username, request.Password)
	if err != nil {
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
