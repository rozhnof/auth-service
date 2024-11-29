package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rozhnof/auth-service/internal/application/services"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/oauth2"
)

type GoogleAuthHandler struct {
	log         *slog.Logger
	authService *services.AuthService
	tracer      trace.Tracer
	cfg         oauth2.Config
}

func NewGoogleAuthHandler(cfg oauth2.Config, service *services.AuthService, log *slog.Logger, tracer trace.Tracer) *GoogleAuthHandler {
	return &GoogleAuthHandler{
		cfg:         cfg,
		authService: service,
		log:         log,
		tracer:      tracer,
	}
}

// @Summary Google OAuth Login
// @Description Redirects to Google OAuth login page.
// @Tags auth
// @Success 303 {object} string "Redirecting to Google OAuth"
// @Router /auth/google/login [get]
func (h *GoogleAuthHandler) Login(c *gin.Context) {
	authURL := h.cfg.AuthCodeURL("randomstate")

	c.Redirect(http.StatusSeeOther, authURL)
}

func (h *GoogleAuthHandler) Callback(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "GoogleAuthHandler.Callback")
	defer span.End()

	state := c.Query("state")
	if state != "randomstate" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "states don't match"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required parameter: code"})
		return
	}

	token, err := h.cfg.Exchange(ctx, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not exchange code for token"})
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get user info"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user data"})
		return
	}

	userDataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read user data"})
		return
	}

	var userData GoogleUserData
	if err := json.Unmarshal(userDataBytes, &userData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unmarshal user data"})
		return
	}

	user, err := h.authService.OAuthLogin(ctx, userData.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to authenticate user"})
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
