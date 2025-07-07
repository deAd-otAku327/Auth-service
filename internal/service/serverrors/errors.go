package serverrors

import "errors"

// Service errors.
var (
	ErrUserGUIDInvalid       = errors.New("invalid user guid")
	ErrIpAddressInvalid      = errors.New("invalid ip address fromat")
	ErrAccessTokenGeneration = errors.New("access token generation failed")
	ErrHashingProcess        = errors.New("hashing process failed")
	ErrSessionAlreadyExists  = errors.New("active session is already exist")
	ErrGUIDExtractionFailed  = errors.New("user guid extraction from context failed")
	ErrGetSession            = errors.New("get session failed")
	ErrCreateSession         = errors.New("create session failed")
	ErrDeleteSession         = errors.New("delete session failed")
)
