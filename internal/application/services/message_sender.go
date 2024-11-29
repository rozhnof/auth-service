package services

import (
	"context"
)

type LoginMessage struct {
	Email string `json:"email"`
}

type RegisterMessage struct {
	Email string `json:"email"`
}

type LoginMessageSender interface {
	SendMessage(context.Context, LoginMessage) error
}

type RegisterMessageSender interface {
	SendMessage(context.Context, RegisterMessage) error
}
