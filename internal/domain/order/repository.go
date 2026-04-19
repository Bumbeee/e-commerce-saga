package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Status string

const (
	StatusCreated   Status = "created"
	StatusPaid      Status = "paid"
	StatusShipped   Status = "shipped"
	StatusDelivered Status = "delivered"
	StatusCancelled Status = "cancelled"
)

type Order struct {
	ID          uuid.UUID       `json:"id,omitempty"`
	UserID      uuid.UUID       `json:"user_id,omitempty"`
	Status      Status          `json:"status"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	CreatedAt   time.Time       `json:"created_at,omitempty"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty"`
}

type Repository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status, expectedUpdatedAt *time.Time) error

	List(ctx context.Context, params ListParams) ([]*Order, int64, error)
}
