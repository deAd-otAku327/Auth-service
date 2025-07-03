package postgres

import (
	"auth-service/internal/config"
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

type postgresDB struct {
	db     *sql.DB
	logger *slog.Logger
}

func New(cfg config.DBConn, logger *slog.Logger) (*postgresDB, error) {
	database, err := sql.Open(PQDriverName, cfg.URL)
	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(cfg.MaxOpenConns)

	return &postgresDB{
		db:     database,
		logger: logger,
	}, nil
}
