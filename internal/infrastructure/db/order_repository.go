package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"e-commerce/internal/domain/order"
)

type pgxRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) order.Repository {
	return &pgxRepository{pool: pool}
}

func (r *pgxRepository) Create(ctx context.Context, ord *order.Order) error {
	query := `
		INSERT INTO orders (user_id, status, total_amount) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at
	`
	err := r.pool.QueryRow(ctx, query, ord.UserID, ord.Status, ord.TotalAmount).
		Scan(&ord.ID, &ord.CreatedAt)

	if err != nil {
		return fmt.Errorf("repository.Create: %w", err)
	}
	return nil
}

func (r *pgxRepository) GetByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	query := `
		SELECT id, user_id, status, total_amount, created_at, updated_at 
		FROM orders WHERE id = $1
	`
	ord := &order.Order{}
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&ord.ID, &ord.UserID, &ord.Status, &ord.TotalAmount, &ord.CreatedAt, &ord.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, order.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("repository.GetById: %w", err)
	}
	return ord, nil
}

func (r *pgxRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = now() WHERE id = $2`
	cmd, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("repository.UpdateStatus: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return order.ErrNotFound
	}
	return nil
}

func (r *pgxRepository) List(ctx context.Context, params order.ListParams) ([]*order.Order, int64, error) {
	params.WithDefaults()
	var (
		query      string
		orders     []*order.Order
		args       []any
		totalCount int64
	)

	if params.UserID != uuid.Nil {
		query = `SELECT id, user_id, status, total_amount, created_at, updated_at, COUNT(*) OVER() AS total_count
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`
		args = append(args, params.UserID, params.Limit, params.Offset)
	} else {
		query = `SELECT id, user_id, status, total_amount, created_at, updated_at, COUNT(*) OVER() AS total_count
		FROM orders
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`
		args = append(args, params.Limit, params.Offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("repository.List: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		ord := &order.Order{}
		var rowTotalCount int64
		if err := rows.Scan(
			&ord.ID,
			&ord.UserID,
			&ord.Status,
			&ord.TotalAmount,
			&ord.CreatedAt,
			&ord.UpdatedAt,
			&rowTotalCount,
		); err != nil {
			return nil, 0, fmt.Errorf("repository.List: %w", err)
		}
		if totalCount == 0 {
			totalCount = rowTotalCount
		}
		orders = append(orders, ord)

	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("repository.List: %w", err)
	}

	return orders, totalCount, nil
}
