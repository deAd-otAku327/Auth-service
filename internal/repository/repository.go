package repository

import (
	"auth-service/internal/config"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/types/models"
	"context"
	"log/slog"
)

type Repository interface {
	GetSession(ctx context.Context, token string) (*models.Session, error)
	CreateSession(ctx context.Context, session *models.Session) error
	RevokeSession(ctx context.Context, token string) error
}

func NewPostgresRepo(cfg config.DBConn, logger *slog.Logger) (Repository, error) {
	return postgres.New(cfg, logger)
}
