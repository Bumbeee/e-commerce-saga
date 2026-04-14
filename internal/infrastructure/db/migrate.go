package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var migrations = []string{
	`CREATE TABLE IF NOT EXISTS orders (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'created',
		total_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		updated_at TIMESTAMPTZ
	)`,
	`CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id)`,
	`CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status)`,
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	for i, query := range migrations {
		if _, err := pool.Exec(ctx, query); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}
	return nil
}
