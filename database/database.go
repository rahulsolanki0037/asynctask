package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, databaseUrl string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, databaseUrl)
}
