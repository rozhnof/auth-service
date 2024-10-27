package handlers

import (
	http_user_mappers "auth/internal/presentation/http/user/dto/mappers"
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
	Email string `json:"email"`
}

type RegisterResponse struct {
	User User `json:"user"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	registeredUser, err := h.service.Register(c.Request.Context(), request.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	response := RegisterResponse{
		User: http_user_mappers.UserToDTO(*registeredUser),
	}

	c.JSON(http.StatusOK, response)
}
