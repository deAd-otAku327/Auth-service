package service

import (
	"auth-service/internal/mappers/dtomap"
	"auth-service/internal/repository"
	"auth-service/internal/service/serverrors"
	"auth-service/internal/tokenizer"
	"auth-service/internal/types/dto"
	"auth-service/internal/types/models"
	"auth-service/internal/types/queries"
	"auth-service/pkg/cryptor"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Login(ctx context.Context, login *models.Login) (*dto.LoginResponse, *http.Cookie, error)
	GetCurrentUser(ctx context.Context) (*dto.UserResponse, error)
	Refresh(ctx context.Context) (*dto.RefreshResponse, error)
	Logout(ctx context.Context) error
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

func (s *authService) Login(ctx context.Context, login *models.Login) (*dto.LoginResponse, *http.Cookie, error) {
	err := uuid.Validate(login.UserGUID)
	if err != nil {
		return nil, nil, fmt.Errorf("%w : %w", serverrors.ErrUserGUIDInvalid, err)
	}

	if len(login.IP) > 15 {
		return nil, nil, fmt.Errorf("%w : %s", serverrors.ErrIpAddressInvalid, login.IP)
	}

	session, err := s.repo.GetSession(ctx, &queries.GetSessionQuery{
		UserGUID:  login.UserGUID,
		UserAgent: login.UserAgent,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("%w : %w", serverrors.ErrGetSession, err)
	}

	if session != nil {
		if time.Now().After(session.ExpiresAt) {
			err := s.repo.DeleteSession(ctx, session.ID)
			if err != nil {
				return nil, nil, fmt.Errorf("%w : %w", serverrors.ErrDeleteSession, err)
			}
		}
		return nil, nil, fmt.Errorf("%w : %w", serverrors.ErrSessionAlreadyExists, err)
	}

	accessToken, err := s.tokenizer.GenerateAccessTokenJWT(login.UserGUID)
	if err != nil {
		return nil, nil, fmt.Errorf("%w : %w", serverrors.ErrAccessTokenGeneration, err)
	}

	refreshCookie := s.tokenizer.GenerateRefreshTokenCookie()

	refreshTokenHash, err := s.cryptor.EncryptKeyword(refreshCookie.Value)
	if err != nil {
		return nil, nil, fmt.Errorf("%w : %w", serverrors.ErrHashingProcess, err)
	}

	err = s.repo.CreateSession(ctx, &queries.CreateSessionQuery{
		UserGUID:     login.UserGUID,
		RefreshToken: refreshTokenHash,
		UserAgent:    login.UserAgent,
		IP:           login.IP,
		ExpiresAt:    refreshCookie.Expires,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("%w : %w", serverrors.ErrCreateSession, err)
	}

	return dtomap.MapToLoginResponse(*accessToken), refreshCookie, nil
}

func (s *authService) GetCurrentUser(ctx context.Context) (*dto.UserResponse, error) {
	return nil, nil
}

func (s *authService) Refresh(ctx context.Context) (*dto.RefreshResponse, error) {
	return nil, nil
}

func (s *authService) Logout(ctx context.Context) error {
	return nil
}
