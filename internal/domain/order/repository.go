package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Order struct {
	Id          uuid.UUID       `json:"id,omitempty"`
	UserId      uuid.UUID       `json:"user_id,omitempty"`
	Status      string          `json:"status"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	CreatedAt   time.Time       `json:"created_at,omitempty"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty"`
}

type Repository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	UpdateStatus(ctx context.Context, id, status string) error
}
