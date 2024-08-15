package handlers

import (
	"auth/internal/application/handlers/http/dto"
	"auth/internal/domain/services"
	repository "auth/internal/infrastructure"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Login @Summary User login
// @Description Logs in a user using their user_id query parameter and returns an access token and refresh token.
// @ID login-user
// @Accept  json
// @Produce  json
// @Param user_id query string true "User ID"
// @Success 200 {object} dto.LoginResponse "Successful login"
// @Failure 400 {string} string "Bad Request - missing or invalid user_id"
// @Failure 404 {string} string "Not Found - user not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has(queryParamUserID) {
		http.Error(w, "query param user_id is required", http.StatusBadRequest)
		return
	}

	userIDStr := r.URL.Query().Get(queryParamUserID)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("query param user_id is invalid: %s", err.Error()), http.StatusBadRequest)
		return
	}
	userIP := r.RemoteAddr

	tokenPair, err := h.service.Login(r.Context(), userID, userIP)
	if err != nil {
		if errors.Is(err, services.ErrUnauthorizedRefresh) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if errors.Is(err, repository.ErrNoData) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "login failed", http.StatusInternalServerError)
		return
	}

	resp := dto.LoginResponse{
		AccessToken:  tokenPair.AccessToken.Token,
		RefreshToken: tokenPair.RefreshToken.Token,
	}
	respJson, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}
