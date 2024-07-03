package repo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	l "github.com/kuromii5/time-tracker/pkg/logger"
)

type DB struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func New(dbURL string, log *slog.Logger) (*DB, error) {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Error("failed to parse config", l.Err(err))

		return nil, fmt.Errorf("%s: %w", "repo.New", err)
	}

	// pool settings
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 5 * time.Minute
	config.MaxConnLifetime = 1 * time.Hour
	log.Debug("connection pool settings",
		slog.Int("max_conns", int(config.MaxConns)),
		slog.Int("min_conns", int(config.MinConns)),
		slog.Duration("max_conn_idle_time", config.MaxConnIdleTime),
		slog.Duration("max_conn_lifetime", config.MaxConnLifetime),
	)

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Error("failed to create connection pool", l.Err(err))

		return nil, fmt.Errorf("%s: %w", "repo.New", err)
	}

	log.Debug("database connection pool created")
	return &DB{pool: pool, log: log}, nil
}

func (db *DB) Close() {
	db.log.Info("closing db connection")

	db.pool.Close()
}
