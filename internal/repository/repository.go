package repository

import (
	"auth-service/internal/config"
	"auth-service/internal/repository/postgres"
	"log/slog"
)

type Repository interface {
}

func NewPostgresRepo(cfg config.DBConn, logger *slog.Logger) (Repository, error) {
	return postgres.New(cfg, logger)
}
