package serverrors

import "errors"

// Service errors.
var (
	ErrUserGUIDInvalid       = errors.New("invalid user guid")
	ErrIpAddressInvalid      = errors.New("invalid ip address fromat")
	ErrAccessTokenGeneration = errors.New("access token generation failed")
	ErrHashingProcess        = errors.New("hashing process failed")
	ErrSessionAlreadyExists  = errors.New("active session is already exist")
	ErrGUIDExtraction        = errors.New("user guid extraction from context failed")
	ErrOldAccessTokenInvalid = errors.New("old access token invalid")
	ErrNoRefreshSession      = errors.New("no refresh session found")
	ErrTokenPairInvalid      = errors.New("invalid token pair")
	ErrRefreshTokenInvalid   = errors.New("invalid refresh token")
	ErrGetSession            = errors.New("get session failed")
	ErrCreateSession         = errors.New("create session failed")
	ErrDeleteSession         = errors.New("delete session failed")
	ErrRenewSession          = errors.New("renew session failed")
)
