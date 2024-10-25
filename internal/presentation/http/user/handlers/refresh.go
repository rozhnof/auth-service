package http_user_handlers

import (
	http_user_requests "auth/internal/presentation/http/user/dto/requests"
	http_user_responses "auth/internal/presentation/http/user/dto/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Refresh @Summary Token refresh
// @Description Refreshes the access token using a valid refresh token.
// @ID refresh-token
// @Accept  json
// @Produce  json
// @Param requestBody body http_user_requests.RefreshRequest true "Refresh token details"
// @Success 200 {object} http_user_responses.RefreshResponse "Tokens refreshed successfully"
// @Failure 400 {string} string "Bad Request - invalid input"
// @Failure 401 {string} string "Unauthorized - invalid refresh token"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var request http_user_requests.RefreshRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var (
		ip           string
		accessToken  string
		refreshToken string
	)

	if err := h.service.Refresh(c.Request.Context(), ip, accessToken, refreshToken); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var response http_user_responses.RefreshResponse

	c.JSON(http.StatusOK, response)
}
