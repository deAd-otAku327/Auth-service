package apierrors

import "errors"

// API errors.
var (
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrAuthorizeNotNeeded   = errors.New("authorization is not needed")
	ErrNoRefreshToken       = errors.New("no refresh token")
	ErrNoAccessToken        = errors.New("no access token to refresh")
	ErrRefreshUnavalible    = errors.New("refresh unavailible")
	ErrInvalidRequestFormat = errors.New("invalid request format")
	ErrInvalidRequestData   = errors.New("invalid request data")
	ErrSomethingWentWrong   = errors.New("sorry, something went wrong")
)
