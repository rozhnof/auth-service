package services

import "context"

type RegisterMessage struct {
	Email       string `json:"email"`
	ConfirmLink string `json:"confirm_link"`
}

type RegisterMessageSender interface {
	SendMessage(context.Context, RegisterMessage) error
}
