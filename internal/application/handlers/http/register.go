package handlers

import (
	"auth/internal/application/handlers/http/dto"
	"auth/pkg/validator"
	"encoding/json"
	"io"
	"net/http"
)

// Register @Summary Register a new user
// @Description Registers a new user by providing an email address and other necessary information.
// @ID register-user
// @Accept  json
// @Produce  json
// @Param requestBody body dto.RegisterRequest true "User registration details"
// @Success 200 {object} dto.RegisterResponse "Successfully registered new user"
// @Failure 400 {string} string "Bad Request - invalid input"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req dto.RegisterRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validator.ValidJsonFields(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.RegisterResponse{
		User: dto.UserToDTO(*user),
	}
	respJson, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}
