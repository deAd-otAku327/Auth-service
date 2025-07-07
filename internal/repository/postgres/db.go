package postgres

import (
	"auth-service/internal/config"
	"auth-service/internal/types/models"
	"auth-service/internal/types/queries"
	"context"
	"database/sql"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
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

func (s *postgresDB) GetSession(ctx context.Context, getSession *queries.GetSessionQuery) (*models.Session, error) {
	query, args, err := sq.Select("*").
		From(SessionsTable).
		Where(sq.Eq{UserIDColumn: getSession.UserGUID, UserAgentColumn: getSession.UserAgent}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRowContext(ctx, query, args...)

	var session models.Session

	err = row.Scan(&session.ID, &session.UserGUID, &session.RefreshToken,
		&session.UserAgent, &session.IP, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}

func (s *postgresDB) GetSessionByToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	query, args, err := sq.Select("*").
		From(SessionsTable).
		Where(sq.Eq{RefreshTokenColumn: refreshToken}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRowContext(ctx, query, args...)

	var session models.Session

	err = row.Scan(&session.ID, &session.UserGUID, &session.RefreshToken,
		&session.UserAgent, &session.IP, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}

func (s *postgresDB) CreateSession(ctx context.Context, createSession *queries.CreateSessionQuery) error {
	query, args, err := sq.Insert(SessionsTable).
		Columns(UserIDColumn, RefreshTokenColumn, UserAgentColumn, IPColumn, ExpiresAtColumn).
		Values(createSession.UserGUID, createSession.RefreshToken, createSession.UserAgent, createSession.IP, createSession.ExpiresAt).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s *postgresDB) DeleteSession(ctx context.Context, sessionID string) error {
	query, args, err := sq.Delete(SessionsTable).
		Where(sq.Eq{SessionIDColumn: sessionID}).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
