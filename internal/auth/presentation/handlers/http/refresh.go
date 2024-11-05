package http_handlers

import (
	"auth/internal/auth/application/services"
	"log/slog"
	"net/http"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"

	"github.com/gin-gonic/gin"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Refresh @Summary Refresh access token
// @Description Refreshes the access token using the refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body RefreshRequest true "Refresh Request"
// @Success 200 {object} RefreshResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	ctx, span := otel.Tracer("tracerName").Start(c.Request.Context(), "Refresh User")
	defer span.End()

	log := h.log.With(
		slog.String("function", "AuthHandler.Refresh"),
	)

	var request RefreshRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Info("bad request")

		c.String(http.StatusBadRequest, err.Error())
		return
	}

	at, rt, err := h.userService.Refresh(ctx, request.RefreshToken)
	if err != nil {
		log.Info("refresh failed", slog.String("error", err.Error()))

		if errors.Is(err, services.ErrUnauthorizedRefresh) {
			c.String(http.StatusUnauthorized, "invalid refresh token")
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	response := RefreshResponse{
		AccessToken:  at,
		RefreshToken: rt,
	}

	c.JSON(http.StatusOK, response)
}
