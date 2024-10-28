package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register @Summary Register a new user
// @Description Registers a new user by providing an email address and other necessary information.
// @ID register-user
// @Accept  json
// @Produce  json
// @Param requestBody body RegisterRequest true "User registration details"
// @Success 200 {object} RegisterResponse "Successfully registered new user"
// @Failure 400 {string} string "Bad Request - invalid input"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/register [post]

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	User
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	registeredUser, err := h.userService.Register(c.Request.Context(), request.Username, request.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	response := RegisterResponse{
		User: UserToDTO(*registeredUser),
	}

	c.JSON(http.StatusOK, response)
}
