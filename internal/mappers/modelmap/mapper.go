package modelmap

import (
	"auth-service/internal/types/dto"
	"auth-service/internal/types/models"
)

func MapToLoginModel(request *dto.LoginRequest, ua, ip string) *models.Login {
	return &models.Login{
		UserGUID:  request.UserGUID,
		UserAgent: ua,
		IP:        ip,
	}
}
