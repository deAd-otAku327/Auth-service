package service

import (
	"auth-service/internal/repository"
	"auth-service/internal/types/dto"
	"auth-service/internal/types/models"
	"context"
	"log/slog"
)

type Service interface {
	Login(ctx context.Context, login *models.Login) (*dto.LoginResponse, *dto.ErrorResponse)
	GetCurrentUser(ctx context.Context) (*dto.UserResponse, *dto.ErrorResponse)
	Refresh(ctx context.Context) (*dto.RefreshResponse, *dto.ErrorResponse)
	Logout(ctx context.Context) *dto.ErrorResponse
}

type authService struct {
	repo   repository.Repository
	logger *slog.Logger
}

func New(repo repository.Repository, logger *slog.Logger) Service {
	return &authService{
		repo:   repo,
		logger: logger,
	}
}

func (s *authService) Login(ctx context.Context, login *models.Login) (*dto.LoginResponse, *dto.ErrorResponse) {
	return nil, nil
}

func (s *authService) GetCurrentUser(ctx context.Context) (*dto.UserResponse, *dto.ErrorResponse) {
	return nil, nil
}

func (s *authService) Refresh(ctx context.Context) (*dto.RefreshResponse, *dto.ErrorResponse) {
	return nil, nil
}

func (s *authService) Logout(ctx context.Context) *dto.ErrorResponse {
	return nil
}
