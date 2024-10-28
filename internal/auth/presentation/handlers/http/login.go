package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login @Summary User login
// @Description Logs in a user using their user_id query parameter and returns an access token and refresh token.
// @ID login-user
// @Accept  json
// @Produce  json
// @Param user_id query string true "User ID"
// @Param requestBody body LoginRequest true "User login details"
// @Success 200 {object} LoginResponse "Successful login"
// @Failure 400 {string} string "Bad Request - missing or invalid user_id"
// @Failure 404 {string} string "Not Found - user not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/login [post]

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	at, rt, err := h.userService.Login(c.Request.Context(), request.Username, request.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	response := LoginResponse{
		AccessToken:  at,
		RefreshToken: rt,
	}

	c.JSON(http.StatusOK, response)
}
