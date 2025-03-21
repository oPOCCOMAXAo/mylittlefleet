package models

import "errors"

var (
	ErrFlowBroken  = errors.New("flow broken")
	ErrInvalidAuth = errors.New("invalid auth")
	ErrNotFound    = errors.New("not found")
)
