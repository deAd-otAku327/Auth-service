package apierrors

import "errors"

// API errors.
var (
	ErrAuthorizeNotNeeded   = errors.New("authorization is not needed")
	ErrInvalidRequestFormat = errors.New("invalid request format")
	ErrInvalidRequestData   = errors.New("invalid request data")
	ErrSomethingWentWrong   = errors.New("sorry, something went wrong")
)
