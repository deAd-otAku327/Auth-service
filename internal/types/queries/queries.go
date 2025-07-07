package queries

import "time"

type GetSessionQuery struct {
	UserGUID  string
	UserAgent string
}

type CreateSessionQuery struct {
	UserGUID     string
	RefreshToken string
	UserAgent    string
	IP           string
	PairID       string
	ExpiresAt    time.Time
}
