package handlers

import (
	"auth/internal/application/handlers/http/dto"
	"auth/internal/domain/services"
	"auth/pkg/validator"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// Refresh @Summary Token refresh
// @Description Refreshes the access token using a valid refresh token.
// @ID refresh-token
// @Accept  json
// @Produce  json
// @Param requestBody body dto.RefreshRequest true "Refresh token details"
// @Success 200 {object} dto.RefreshResponse "Tokens refreshed successfully"
// @Failure 400 {string} string "Bad Request - invalid input"
// @Failure 401 {string} string "Unauthorized - invalid refresh token"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req dto.RefreshRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.ValidJsonFields(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userIP := r.RemoteAddr

	tokenPair, err := h.service.Refresh(r.Context(), userIP, req.AccessToken, req.RefreshToken)
	if errors.Is(err, services.ErrUnauthorizedRefresh) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.RefreshResponse{
		AccessToken:  tokenPair.AccessToken.Token,
		RefreshToken: tokenPair.RefreshToken.Token,
	}
	respJson, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}
