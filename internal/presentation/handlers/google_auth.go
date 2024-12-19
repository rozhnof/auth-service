package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rozhnof/auth-service/internal/application/services"
	"github.com/rozhnof/auth-service/internal/presentation/clients"
	"go.opentelemetry.io/otel/trace"
)

const (
	todoState = "randomstate"
)

type GoogleAuthHandler struct {
	authClient  *clients.GoogleAuthClient
	log         *slog.Logger
	authService *services.AuthService
	tracer      trace.Tracer
}

func NewGoogleAuthHandler(authClient *clients.GoogleAuthClient, service *services.AuthService, log *slog.Logger, tracer trace.Tracer) *GoogleAuthHandler {
	return &GoogleAuthHandler{
		authClient:  authClient,
		authService: service,
		log:         log,
		tracer:      tracer,
	}
}

// @Summary Google OAuth Login
// @Description Redirects to Google OAuth login page.
// @Tags Auth
// @Success 303 {object} string "Redirecting to Google OAuth"
// @Router /auth/google/login [get]
func (h *GoogleAuthHandler) Login(c *gin.Context) {
	authURL := h.authClient.GetAuthURL(todoState)

	c.Redirect(http.StatusSeeOther, authURL)
}

type callbackQueryParams struct {
	State string `form:"state" binding:"required"`
	Code  string `form:"code" binding:"required"`
}

func (h *GoogleAuthHandler) Callback(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "GoogleAuthHandler.Callback")
	defer span.End()

	var queryParams callbackQueryParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if queryParams.State != todoState {
		c.String(http.StatusBadRequest, "states don't match")
		return
	}

	token, err := h.authClient.ExchangeOAuthToken(ctx, queryParams.Code)
	if err != nil {
		c.String(http.StatusInternalServerError, "could not exchange code for token")
		return
	}

	userInfo, statusCode, err := h.authClient.GetUserInfo(token.AccessToken)
	if err != nil {
		h.log.Warn("failed get user info", slog.String("error", err.Error()))

		c.Status(http.StatusInternalServerError)
		return
	}

	if statusCode < 200 || statusCode >= 400 {
		c.Status(statusCode)
		return
	}

	user, err := h.authService.OAuthLogin(ctx, userInfo.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to authenticate user")
		return
	}

	response := GoogleLoginResponse{
		UserID:       user.ID(),
		Email:        user.Email(),
		AccessToken:  user.AccessToken().Token(),
		RefreshToken: user.RefreshToken().Token(),
	}

	c.JSON(http.StatusOK, response)
}
