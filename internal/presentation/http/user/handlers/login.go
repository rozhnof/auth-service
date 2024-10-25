package http_user_handlers

import (
	repository "auth/internal/infrastructure"
	http_user_requests "auth/internal/presentation/http/user/dto/requests"
	http_user_responses "auth/internal/presentation/http/user/dto/responses"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Login @Summary User login
// @Description Logs in a user using their user_id query parameter and returns an access token and refresh token.
// @ID login-user
// @Accept  json
// @Produce  json
// @Param user_id query string true "User ID"
// @Param requestBody body http_user_requests.LoginRequest true "User login details"
// @Success 200 {object} http_user_responses.LoginResponse "Successful login"
// @Failure 400 {string} string "Bad Request - missing or invalid user_id"
// @Failure 404 {string} string "Not Found - user not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var request http_user_requests.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var userID uuid.UUID
	if err := h.service.Login(c.Request.Context(), userID); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var response http_user_responses.LoginResponse

	c.JSON(http.StatusOK, response)
}
