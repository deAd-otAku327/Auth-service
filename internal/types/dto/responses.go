package dto

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type RefreshResponse struct {
	NewAccessToken string `json:"access_token"`
}

type UserResponse struct {
	UserGUID string `json:"user_guid"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
