package service

import (
	"auth-service/internal/repository"
	"log/slog"
)

type Service interface {
}

type authService struct {
	repo   repository.Repository
	logger *slog.Logger
}

func New(repo repository.Repository, logger *slog.Logger) Service {
	return authService{
		repo:   repo,
		logger: logger,
	}
}
