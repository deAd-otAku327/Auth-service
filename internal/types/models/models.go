package models

import "time"

type Login struct {
	UserGUID  string
	UserAgent string
	IP        string
}

type Refresh struct {
	AccessToken  string
	RefreshToken string
}

type User struct {
	GUID string
}

type Session struct {
	ID           string
	UserGUID     string
	RefreshToken string
	UserAgent    string
	IP           string
	PairID       string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}
