package models

import "errors"

var (
	ErrInvalidAuth = errors.New("invalid auth")
	ErrNotFound    = errors.New("not found")
)
