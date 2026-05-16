package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Conn(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, err
}
