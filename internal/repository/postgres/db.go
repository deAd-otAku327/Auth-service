package postgres

import (
	"auth-service/internal/config"
	"auth-service/internal/types/models"
	"auth-service/internal/types/queries"
	"context"
	"database/sql"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

type postgresDB struct {
	db     *sql.DB
	logger *slog.Logger
}

func New(cfg config.DBConn, logger *slog.Logger) (*postgresDB, error) {
	database, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(cfg.MaxOpenConns)

	migrationsDir := "file://../migrations"

	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	migrator, err := migrate.NewWithDatabaseInstance(migrationsDir, "auth-db", driver)
	if err != nil {
		return nil, err
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	slog.Info("database migrated and ready")

	return &postgresDB{
		db:     database,
		logger: logger,
	}, nil
}

func (db *postgresDB) GetSession(ctx context.Context, getSession *queries.GetSessionQuery) (*models.Session, error) {
	return nil, nil
}

func (db *postgresDB) GetSessionByToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	return nil, nil
}

func (db *postgresDB) CreateSession(ctx context.Context, createSession *queries.CreateSessionQuery) error {
	return nil
}

func (db *postgresDB) DeleteSession(ctx context.Context, sessionID string) error {
	return nil
}
