package repository

import "errors"

var (
	ErrNotExists = errors.New("not exists")
	ErrDuplicate = errors.New("duplicate")
)
