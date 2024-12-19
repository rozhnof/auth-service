package handlers

import "github.com/google/uuid"

type GoogleLoginResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	Email        string    `json:"email"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}
