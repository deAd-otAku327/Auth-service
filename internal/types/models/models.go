package models

import "time"

type User struct {
	GUID string
}

type RefreshToken struct {
	UserGUID  string
	Token     string
	UserAgent string
	IP        string
	ExpiresAt time.Time
	CreatedAt time.Time
}
