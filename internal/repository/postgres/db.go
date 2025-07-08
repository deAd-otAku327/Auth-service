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

func (s *postgresDB) GetSession(ctx context.Context, getSessionQuery *queries.GetSessionQuery) (*models.Session, error) {
	query, args, err := sq.Select("*").
		From(SessionsTable).
		Where(sq.Eq{UserIDColumn: getSessionQuery.UserGUID, UserAgentColumn: getSessionQuery.UserAgent}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", repoerrors.ErrQueryBuilding, err)
	}

	row := s.db.QueryRowContext(ctx, query, args...)

	var session models.Session

	err = row.Scan(&session.ID, &session.UserGUID, &session.RefreshToken,
		&session.UserAgent, &session.IP, &session.PairID, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("%w: %w", repoerrors.ErrQueryExec, err)
	}

	return &session, nil
}

func (s *postgresDB) CreateSession(ctx context.Context, createSessionquery *queries.CreateSessionQuery) error {
	err := createSession(ctx, s.db, createSessionquery)
	if err != nil {
		return err
	}

	return nil
}

func (s *postgresDB) DeleteSession(ctx context.Context, sessionID string) error {
	err := deleteSession(ctx, s.db, sessionID)
	if err != nil {
		return err
	}

	return nil
}

func (s *postgresDB) RenewSession(ctx context.Context, oldSessionID string, createSessionQuery *queries.CreateSessionQuery) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%w: %w", repoerrors.ErrTransactionBegin, err)
	}

	err = deleteSession(ctx, tx, oldSessionID)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return fmt.Errorf("%w + %w", err, txErr)
		}
		return err
	}

	err = createSession(ctx, tx, createSessionQuery)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return fmt.Errorf("%w + %w", err, txErr)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%w: %w", repoerrors.ErrTransactionCommit, err)
	}

	return nil
}
