package dto

type LoginRequest struct {
	UserGUID string
}

type RefreshRequest struct {
	RefreshToken string
}
