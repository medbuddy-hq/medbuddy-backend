package constant

import "errors"

const (
	ErrValidation   = "validation error"
	ErrServer       = "server error"
	ErrUnauthorized = "unauthorized"
	ErrBinding      = "binding error"
	ErrRequest      = "could not execute request"
)

var (
	ErrResourceAlreadyExists = errors.New("resource already exists")
)
