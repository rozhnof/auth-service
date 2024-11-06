//go:generate ffjson $GOFILE
package http_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	User
}

// Register @Summary User registration
// @Description Registers a new user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param register body RegisterRequest true "Register Request"
// @Success 200 {object} RegisterResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "AuthHandler.Register")
	defer span.End()

	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	registeredUser, err := h.userService.Register(ctx, request.Username, request.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	response := RegisterResponse{
		User: UserToDTO(*registeredUser),
	}

	bytes, _ := response.MarshalJSON()

	c.Data(http.StatusOK, "application/json", bytes)
}
