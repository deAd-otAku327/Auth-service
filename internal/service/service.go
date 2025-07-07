package service

import (
	"auth-service/internal/apperrors"
	"auth-service/internal/mappers/dtomap"
	"auth-service/internal/repository"
	"auth-service/internal/tokenizer"
	"auth-service/internal/types/dto"
	"auth-service/internal/types/models"
	"auth-service/internal/types/queries"
	"auth-service/pkg/cryptor"
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Service interface {
	Login(ctx context.Context, login *models.Login) (*dto.LoginResponse, *http.Cookie, *dto.ErrorResponse)
	GetCurrentUser(ctx context.Context) (*dto.UserResponse, *dto.ErrorResponse)
	Refresh(ctx context.Context) (*dto.RefreshResponse, *dto.ErrorResponse)
	Logout(ctx context.Context) *dto.ErrorResponse
}

type authService struct {
	repo      repository.Repository
	cryptor   cryptor.Cryptor
	tokenizer tokenizer.Tokenizer
	logger    *slog.Logger
}

func New(repo repository.Repository, cryptor cryptor.Cryptor, tok tokenizer.Tokenizer, logger *slog.Logger) Service {
	return &authService{
		repo:      repo,
		cryptor:   cryptor,
		tokenizer: tok,
		logger:    logger,
	}
}

func (s *authService) Login(ctx context.Context, login *models.Login) (*dto.LoginResponse, *http.Cookie, *dto.ErrorResponse) {
	session, err := s.repo.GetSession(ctx, &queries.GetSessionQuery{
		UserGUID:  login.UserGUID,
		UserAgent: login.UserAgent,
	})
	if err != nil {
		return nil, nil, dtomap.MapToErrorResponse(apperrors.ErrSomethingWentWrong, http.StatusInternalServerError)
	}

	if session != nil {
		if time.Now().After(session.ExpiresAt) {
			err := s.repo.DeleteSession(ctx, session.ID)
			if err != nil {
				return nil, nil, dtomap.MapToErrorResponse(apperrors.ErrSomethingWentWrong, http.StatusInternalServerError)
			}
		}
		return nil, nil, dtomap.MapToErrorResponse(apperrors.ErrSessionAlreadyExists, http.StatusFound)
	}

	accessToken, err := s.tokenizer.GenerateAccessTokenJWT(login.UserGUID)
	if err != nil {
		return nil, nil, dtomap.MapToErrorResponse(apperrors.ErrSomethingWentWrong, http.StatusInternalServerError)
	}

	refreshCookie := s.tokenizer.GenerateRefreshTokenCookie()

	refreshTokenHash, err := s.cryptor.EncryptKeyword(refreshCookie.Value)
	if err != nil {
		return nil, nil, dtomap.MapToErrorResponse(apperrors.ErrSomethingWentWrong, http.StatusInternalServerError)
	}

	err = s.repo.CreateSession(ctx, &queries.CreateSessionQuery{
		UserGUID:     login.UserGUID,
		RefreshToken: refreshTokenHash,
		UserAgent:    login.UserAgent,
		IP:           login.IP,
		ExpiresAt:    refreshCookie.Expires,
	})
	if err != nil {
		return nil, nil, dtomap.MapToErrorResponse(apperrors.ErrSomethingWentWrong, http.StatusInternalServerError)
	}

	return dtomap.MapToLoginResponse(*accessToken), refreshCookie, nil
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
