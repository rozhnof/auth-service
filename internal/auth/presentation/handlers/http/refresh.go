package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Refresh @Summary Token refresh
// @Description Refreshes the access token using a valid refresh token.
// @ID refresh-token
// @Accept  json
// @Produce  json
// @Param requestBody body RefreshRequest true "Refresh token details"
// @Success 200 {object} RefreshResponse "Tokens refreshed successfully"
// @Failure 400 {string} string "Bad Request - invalid input"
// @Failure 401 {string} string "Unauthorized - invalid refresh token"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/refresh [post]

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var request RefreshRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	at, rt, err := h.userService.Refresh(c.Request.Context(), request.RefreshToken)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	response := RefreshResponse{
		AccessToken:  at,
		RefreshToken: rt,
	}

	c.JSON(http.StatusOK, response)
}
