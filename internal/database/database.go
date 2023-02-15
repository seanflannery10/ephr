package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slog"
)

func New(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	defer dbpool.Close()

	slog.Info("database connection pool established")

	return dbpool, nil
}
