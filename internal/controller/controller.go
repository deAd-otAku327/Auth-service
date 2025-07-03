package controller

import (
	"auth-service/internal/service"
	"log/slog"
)

type Controller interface {
}

type authController struct {
	service service.Service
	logger  *slog.Logger
}

func New(service service.Service, logger *slog.Logger) Controller {
	return authController{
		service: service,
		logger:  logger,
	}
}
