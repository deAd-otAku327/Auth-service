package models

import "time"

type User struct {
	GUID string
}

type Session struct {
	UserGUID     string
	RefreshToken string
	UserAgent    string
	IP           string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}
