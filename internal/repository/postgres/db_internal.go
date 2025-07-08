package postgres

import (
	"auth-service/internal/repository/repoerrors"
	"auth-service/internal/types/queries"
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type sqlexecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func createSession(ctx context.Context, exec sqlexecutor, createSessionQuery *queries.CreateSessionQuery) error {
	query, args, err := sq.Insert(SessionsTable).
		Columns(UserIDColumn, RefreshTokenColumn, UserAgentColumn, IPColumn, PairIDColumn, ExpiresAtColumn).
		Values(
			createSessionQuery.UserGUID, createSessionQuery.RefreshToken, createSessionQuery.UserAgent,
			createSessionQuery.IP, createSessionQuery.PairID, createSessionQuery.ExpiresAt).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("%w: %w", repoerrors.ErrQueryBuilding, err)
	}

	_, err = exec.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %w", repoerrors.ErrQueryExec, err)
	}

	return nil
}

func deleteSession(ctx context.Context, exec sqlexecutor, sessionID string) error {
	query, args, err := sq.Delete(SessionsTable).
		Where(sq.Eq{SessionIDColumn: sessionID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("%w: %w", repoerrors.ErrQueryBuilding, err)
	}

	_, err = exec.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %w", repoerrors.ErrQueryExec, err)
	}

	return nil
}
