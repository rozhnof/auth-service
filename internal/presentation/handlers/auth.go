package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
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
		if errors.Is(err, services.ErrDuplicate) {
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
		if errors.Is(err, services.ErrObjectNotFound) {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		if errors.Is(err, services.ErrInvalidPassword) {
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
		if errors.Is(err, services.ErrObjectNotFound) {
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
