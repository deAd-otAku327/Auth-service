package postgres

import (
	"auth-service/internal/config"
	"auth-service/internal/repository/repoerrors"
	"auth-service/internal/types/models"
	"auth-service/internal/types/queries"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

type postgresDB struct {
	db *sql.DB
}

func New(cfg config.DBConn) (*postgresDB, error) {
	database, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(cfg.MaxOpenConns)

	migrationsDir := os.Getenv("MIGRATIONS_DIR")

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
		db: database,
	}, nil
}

func (s *postgresDB) GetSession(ctx context.Context, getSession *queries.GetSessionQuery) (*models.Session, error) {
	query, args, err := sq.Select("*").
		From(SessionsTable).
		Where(sq.Eq{UserIDColumn: getSession.UserGUID, UserAgentColumn: getSession.UserAgent}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", repoerrors.ErrQueryBuildingFailed, err)
	}

	row := s.db.QueryRowContext(ctx, query, args...)

	var session models.Session

	err = row.Scan(&session.ID, &session.UserGUID, &session.RefreshToken,
		&session.UserAgent, &session.IP, &session.PairID, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("%w : %w", repoerrors.ErrQueryExecFailed, err)
	}

	return &session, nil
}

func (s *postgresDB) CreateSession(ctx context.Context, createSession *queries.CreateSessionQuery) error {
	query, args, err := sq.Insert(SessionsTable).
		Columns(UserIDColumn, RefreshTokenColumn, UserAgentColumn, IPColumn, PairIDColumn, ExpiresAtColumn).
		Values(
			createSession.UserGUID, createSession.RefreshToken, createSession.UserAgent,
			createSession.IP, createSession.PairID, createSession.ExpiresAt).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("%w: %w", repoerrors.ErrQueryBuildingFailed, err)
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w : %w", repoerrors.ErrQueryExecFailed, err)
	}

	return nil
}

func (s *postgresDB) DeleteSession(ctx context.Context, sessionID string) error {
	query, args, err := sq.Delete(SessionsTable).
		Where(sq.Eq{SessionIDColumn: sessionID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("%w: %w", repoerrors.ErrQueryBuildingFailed, err)
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w : %w", repoerrors.ErrQueryExecFailed, err)
	}

	return nil
}
