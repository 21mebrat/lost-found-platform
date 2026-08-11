package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	// parse db config
	cfg, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return nil, fmt.Errorf("parse db config:%w", err)
	}
	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnIdleTime = 2 * time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	// actual db pool
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create db pool:%w", err)
	}

	// db context
	pgxContex, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()
	// db ping
	if err := pool.Ping(pgxContex); err != nil {
		pool.Close()
		return nil, fmt.Errorf("png db:%w", err)
	}
	return pool, nil
}
