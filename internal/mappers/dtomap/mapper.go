package dtomap

import "auth-service/internal/types/dto"

func MapToLoginResponse(accessToken string) *dto.LoginResponse {
	return &dto.LoginResponse{
		AccessToken: accessToken,
	}
}

func MapToErrorResponse(err error, code int) *dto.ErrorResponse {
	return &dto.ErrorResponse{
		Message: err.Error(),
		Code:    code,
	}
}
