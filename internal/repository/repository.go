package repository

import (
	"auth-service/internal/config"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/types/models"
	"auth-service/internal/types/queries"
	"context"
)

type Repository interface {
	GetSession(ctx context.Context, getSession *queries.GetSessionQuery) (*models.Session, error)
	CreateSession(ctx context.Context, createSession *queries.CreateSessionQuery) error
	DeleteSession(ctx context.Context, sessionID string) error
	RenewSession(ctx context.Context, oldSessionID string, createSession *queries.CreateSessionQuery) error
}

func NewPostgresRepo(cfg config.DBConn) (Repository, error) {
	return postgres.New(cfg)
}
