package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slog"
)

type Config struct {
	DSN                   string
	MinConns              int32
	MaxConns              int32
	MaxConnLifetime       string
	MaxConnLifetimeJitter string
	MaxConnIdleTime       string
}

func New(cfg Config) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, err
	}

	config.MinConns = cfg.MinConns
	config.MaxConns = cfg.MaxConns

	maxConnLifetime, err := time.ParseDuration(cfg.MaxConnLifetime)
	if err != nil {
		return nil, err
	}

	config.MaxConnLifetime = maxConnLifetime

	maxConnLifetimeJitter, err := time.ParseDuration(cfg.MaxConnLifetimeJitter)
	if err != nil {
		return nil, err
	}

	config.MaxConnLifetimeJitter = maxConnLifetimeJitter

	maxConnIdleTime, err := time.ParseDuration(cfg.MaxConnIdleTime)
	if err != nil {
		return nil, err
	}

	config.MaxConnIdleTime = maxConnIdleTime

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	defer dbpool.Close()

	slog.Info("database connection pool established")

	return dbpool, nil
}
