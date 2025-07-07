package repository

import (
	"auth-service/internal/config"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/types/models"
	"auth-service/internal/types/queries"
	"context"
	"log/slog"
)

type Repository interface {
	GetSession(ctx context.Context, getSession *queries.GetSessionQuery) (*models.Session, error)
	GetSessionByToken(ctx context.Context, refreshToken string) (*models.Session, error)
	CreateSession(ctx context.Context, createSession *queries.CreateSessionQuery) error
	DeleteSession(ctx context.Context, sessionID string) error
}

func NewPostgresRepo(cfg config.DBConn, logger *slog.Logger) (Repository, error) {
	return postgres.New(cfg, logger)
}
