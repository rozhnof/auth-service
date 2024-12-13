package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	repo "github.com/rozhnof/auth-service/internal/application/repository"
	"github.com/rozhnof/auth-service/internal/application/services"
	"go.opentelemetry.io/otel/trace"
)

type AuthHandler struct {
	log         *slog.Logger
	authService *services.AuthService
	tracer      trace.Tracer
}

func NewAuthHandler(service *services.AuthService, log *slog.Logger, tracer trace.Tracer) *AuthHandler {
	return &AuthHandler{
		authService: service,
		log:         log,
		tracer:      tracer,
	}
}

// Confirm godoc
// @Summary Confirm user registration
// @Description This endpoint confirms user registration using the provided email and register_token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param email query string true "User email"
// @Param register_token query string true "Register token"
// @Success 200 {string} string "OK"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/confirm [post]
func (h *AuthHandler) Confirm(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "AuthHandler.Confirm")
	defer span.End()

	var (
		email         = c.Query("email")
		registerToken = c.Query("register_token")
	)

	if err := h.authService.Confirm(ctx, email, registerToken); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// Register @Summary User registration
// @Description Registers a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param register body RegisterRequest true "Register Request"
// @Param refcode query string false "Referral Code"
// @Success 200 {object} RegisterResponse
// @Failure 400 {string} string "Missing required parameters"
// @Failure 409 {string} string "User with this email already exists"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "AuthHandler.Register")
	defer span.End()

	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "missing required parameters")
		return
	}

	registeredUser, err := h.authService.Register(ctx, request.Email, request.Password)
	if err != nil {
		if errors.Is(err, repo.ErrDuplicate) {
			c.String(http.StatusConflict, err.Error())
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	response := RegisterResponse{
		UserID: registeredUser.ID(),
		Email:  registeredUser.Email(),
	}

	c.JSON(http.StatusOK, response)
}

// Login @Summary User login
// @Description Login user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login Request"
// @Success 200 {object} LoginResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "AuthHandler.Login")
	defer span.End()

	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "missing required parameters")
		return
	}

	at, rt, err := h.authService.Login(ctx, request.Email, request.Password)
	if err != nil {
		if errors.Is(err, repo.ErrObjectNotFound) || errors.Is(err, services.ErrInvalidPassword) {
			c.String(http.StatusOK, "invalid email or password")
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		AccessToken:  at,
		RefreshToken: rt,
	}

	c.JSON(http.StatusOK, response)
}

// Refresh @Summary Refresh access token
// @Description Refreshes the access token using the refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param refresh body RefreshRequest true "Refresh Request"
// @Success 200 {object} RefreshResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "AuthHandler.Refresh")
	defer span.End()

	var request RefreshRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "missing required parameters")
		return
	}

	at, rt, err := h.authService.Refresh(ctx, request.RefreshToken)
	if err != nil {
		if errors.Is(err, repo.ErrObjectNotFound) {
			c.String(http.StatusUnauthorized, err.Error())
			return
		}

		if errors.Is(err, services.ErrUnauthorizedRefresh) {
			c.String(http.StatusUnauthorized, err.Error())
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	response := RefreshResponse{
		AccessToken:  at,
		RefreshToken: rt,
	}

	c.JSON(http.StatusOK, response)
}
