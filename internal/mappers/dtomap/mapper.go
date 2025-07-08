package dtomap

import (
	"auth-service/internal/types/dto"
)

func MapToLoginResponse(accessToken string) *dto.LoginResponse {
	return &dto.LoginResponse{
		AccessToken: accessToken,
	}
}

func MapToUserResponse(userGUID string) *dto.UserResponse {
	return &dto.UserResponse{
		UserGUID: userGUID,
	}
}

func MapToErrorResponse(err error, code int) *dto.ErrorResponse {
	return &dto.ErrorResponse{
		Message: err.Error(),
		Code:    code,
	}
}

func MapToRefreshResponse(accessToken string) *dto.RefreshResponse {
	return &dto.RefreshResponse{
		NewAccessToken: accessToken,
	}
}
