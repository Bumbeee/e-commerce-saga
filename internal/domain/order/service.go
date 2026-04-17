package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateOrderDTO struct {
	UserID      string          `json:"user_id"`
	TotalAmount decimal.Decimal `json:"total_amount"`
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOrder(ctx context.Context, dto CreateOrderDTO) (*Order, error) {
	if err := validateCreateDTO(dto); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	userID, err := uuid.Parse(dto.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id format: %w", err)
	}

	order := &Order{
		ID:          uuid.New(),
		UserID:      userID,
		Status:      StatusCreated,
		TotalAmount: dto.TotalAmount,
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("service.CreateOrder: %w", err)
	}

	return order, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Order, error) {
	ord, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service.GetByID: %w", err)
	}
	return ord, nil
}

func validateCreateDTO(dto CreateOrderDTO) error {
	if dto.UserID == "" {
		return errors.New("user_id is required")
	}
	if dto.TotalAmount.LessThanOrEqual(decimal.Zero) {
		return errors.New("total_amount must be positive")
	}
	return nil
}
