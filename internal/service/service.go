package service

import (
	"auth-service/internal/mappers/dtomap"
	"auth-service/internal/middleware"
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
	Refresh(ctx context.Context, refresh *models.Refresh) (*dto.RefreshResponse, *http.Cookie, error)
	Logout(ctx context.Context, userAgent string) error
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
		return nil, nil, fmt.Errorf("%w: %w", serverrors.ErrUserGUIDInvalid, err)
	}

	if len(login.IP) < 7 || len(login.IP) > 15 {
		return nil, nil, fmt.Errorf("%w: %s", serverrors.ErrIpAddressInvalid, login.IP)
	}

	session, err := s.repo.GetSession(ctx, &queries.GetSessionQuery{
		UserGUID:  login.UserGUID,
		UserAgent: login.UserAgent,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", serverrors.ErrGetSession, err)
	}

	if session != nil {
		if time.Now().After(session.ExpiresAt) {
			err := s.repo.DeleteSession(ctx, session.ID)
			if err != nil {
				return nil, nil, fmt.Errorf("%w: %w", serverrors.ErrDeleteSession, err)
			}
		}
		return nil, nil, fmt.Errorf("%w: %w", serverrors.ErrSessionAlreadyExists, err)
	}

	tokenPairID := uuid.NewString()
	accessToken, err := s.tokenizer.GenerateAccessTokenJWT(login.UserGUID, tokenPairID)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", serverrors.ErrAccessTokenGeneration, err)
	}

	refreshCookie := s.tokenizer.GenerateRefreshTokenCookie()

	refreshTokenHash, err := s.cryptor.EncryptKeyword(refreshCookie.Value)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", serverrors.ErrHashingProcess, err)
	}

	err = s.repo.CreateSession(ctx, &queries.CreateSessionQuery{
		UserGUID:     login.UserGUID,
		RefreshToken: refreshTokenHash,
		UserAgent:    login.UserAgent,
		IP:           login.IP,
		PairID:       tokenPairID,
		ExpiresAt:    refreshCookie.Expires,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", serverrors.ErrCreateSession, err)
	}

	return dtomap.MapToLoginResponse(*accessToken), refreshCookie, nil
}

func (s *authService) GetCurrentUser(ctx context.Context) (*dto.UserResponse, error) {
	currUserGUID, ok := ctx.Value(middleware.UserGUIDKey).(string)
	if !ok {
		return nil, serverrors.ErrGUIDExtraction
	}

	return dtomap.MapToUserResponse(currUserGUID), nil
}

func (s *authService) Refresh(ctx context.Context, refresh *models.Refresh) (*dto.RefreshResponse, *http.Cookie, error) {
	if len(refresh.IP) < 7 || len(refresh.IP) > 15 {
		return nil, nil, fmt.Errorf("%w: %s", serverrors.ErrIpAddressInvalid, refresh.IP)
	}

	tokenClaims, err := s.tokenizer.VerifyAccessTokenJWT(refresh.AccessToken, true)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", serverrors.ErrOldAccessTokenInvalid, err)
	}

	userGUID, err := tokenClaims.GetSubject()
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", serverrors.ErrOldAccessTokenInvalid, err)
	}

	accessPairID, ok := tokenClaims[tokenizer.PairClaimsKey].(string)
	if !ok {
		return nil, nil, fmt.Errorf("%w: %s", serverrors.ErrOldAccessTokenInvalid, "no access pair id")
	}

	session, err := s.repo.GetSession(ctx, &queries.GetSessionQuery{
		UserGUID:  userGUID,
		UserAgent: refresh.UserAgent,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", serverrors.ErrGetSession, err)
	}

	if session == nil {
		return nil, nil, serverrors.ErrNoRefreshSession
	}

	if time.Now().After(session.ExpiresAt) || refresh.UserAgent != session.UserAgent {
		err := s.repo.DeleteSession(ctx, session.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %w", serverrors.ErrDeleteSession, err)
		}
		return nil, nil, serverrors.ErrNoRefreshSession
	}

	// TODO: IP checking.

	if accessPairID != session.PairID {
		return nil, nil, serverrors.ErrTokenPairInvalid
	}

	err = s.cryptor.CompareHashAndKeyword(session.RefreshToken, refresh.RefreshToken)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", serverrors.ErrRefreshTokenInvalid, err)
	}

	newTokenPairID := uuid.NewString()
	newAccessToken, err := s.tokenizer.GenerateAccessTokenJWT(session.UserGUID, newTokenPairID)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", serverrors.ErrAccessTokenGeneration, err)
	}

	newRefreshCookie := s.tokenizer.GenerateRefreshTokenCookie()

	newRefreshTokenHash, err := s.cryptor.EncryptKeyword(newRefreshCookie.Value)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", serverrors.ErrHashingProcess, err)
	}

	err = s.repo.RenewSession(ctx, session.ID, &queries.CreateSessionQuery{
		UserGUID:     session.UserGUID,
		RefreshToken: newRefreshTokenHash,
		UserAgent:    refresh.UserAgent,
		IP:           refresh.IP,
		PairID:       newTokenPairID,
		ExpiresAt:    newRefreshCookie.Expires,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", serverrors.ErrRenewSession, err)
	}

	return dtomap.MapToRefreshResponse(*newAccessToken), newRefreshCookie, nil
}

func (s *authService) Logout(ctx context.Context, userAgent string) error {
	currUserGUID, ok := ctx.Value(middleware.UserGUIDKey).(string)
	if !ok {
		return serverrors.ErrGUIDExtraction
	}

	session, err := s.repo.GetSession(ctx, &queries.GetSessionQuery{
		UserGUID:  currUserGUID,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", serverrors.ErrGetSession, err)
	}

	if session != nil {
		err := s.repo.DeleteSession(ctx, session.ID)
		if err != nil {
			return fmt.Errorf("%w: %w", serverrors.ErrDeleteSession, err)
		}
	}

	return nil
}
