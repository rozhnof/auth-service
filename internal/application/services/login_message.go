package services

import "context"

type LoginMessage struct {
	Email string `json:"email"`
}

type LoginMessageSender interface {
	SendMessage(context.Context, LoginMessage) error
}
