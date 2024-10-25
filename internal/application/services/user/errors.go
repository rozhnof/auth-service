package user_services

import "errors"

var (
	ErrUnauthorizedRefresh = errors.New("unauthorized refresh")
	ErrInvalidPassword     = errors.New("invalid password")
)
