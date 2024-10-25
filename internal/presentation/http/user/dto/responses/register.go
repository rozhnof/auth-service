package http_user_responses

import http_user_structs "auth/internal/presentation/http/user/dto/structs"

type RegisterResponse struct {
	User http_user_structs.User `json:"user"`
}
