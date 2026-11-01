package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, conn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("error create connection")
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error ping connection")
	}

	return pool, nil
}
