package repo

import "errors"

var (
	ErrObjectNotFound = errors.New("object not found")
	ErrDuplicate      = errors.New("object is duplicate")
)
