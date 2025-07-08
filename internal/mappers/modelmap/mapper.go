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

func MapToRefreshModel(access, refresh, ua, ip string) *models.Refresh {
	return &models.Refresh{
		AccessToken:  access,
		RefreshToken: refresh,
		UserAgent:    ua,
		IP:           ip,
	}
}
