package domain

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("Internal Server Error")
	// ErrResourceNotFound will throw if the requested item is not exists
	ErrResourceNotFound = errors.New("Unknown resource or resource not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("Conflicting state")
	// ErrBadParamInput will throw if the given request input is not valid
	ErrBadParamInput = errors.New("Bad or invalid input")
)
