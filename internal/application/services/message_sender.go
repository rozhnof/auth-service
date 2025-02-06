package services

import "context"

type MessageSender interface {
	SendMessage(context.Context, any) error
}

type LoginMessage struct {
	Email string `json:"email"`
}

type RegisterMessage struct {
	Email       string `json:"email"`
	ConfirmLink string `json:"confirm_link"`
}
