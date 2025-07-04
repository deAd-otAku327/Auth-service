package dto

type LoginRequest struct {
	UserGUID string `shema:"guid"`
}

type RefreshRequest struct {
	RefreshToken string
}
