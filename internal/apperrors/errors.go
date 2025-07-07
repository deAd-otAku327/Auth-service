package apperrors

import "errors"

// Controller errors.
var (
	ErrInvalidRequestParams = errors.New("invalid request params")
)

// Service errors.
var (
	ErrSomethingWentWrong   = errors.New("sorry, something went wrong")
	ErrSessionAlreadyExists = errors.New("active session is already exist")
)
