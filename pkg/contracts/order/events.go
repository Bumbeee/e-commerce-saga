// pkg/contracts/order/events.go

package order

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// OrderCreatedEvent — контракт события создания заказа.
// Используется для межсервисной коммуникации через Kafka.
type OrderCreatedEvent struct {
	SchemaVersion  string          `json:"schema_version"`
	EventID        uuid.UUID       `json:"event_id"`
	OrderID        uuid.UUID       `json:"order_id"`
	UserID         uuid.UUID       `json:"user_id"`
	TotalAmount    decimal.Decimal `json:"total_amount"`
	CreatedAt      time.Time       `json:"created_at"`
	IdempotencyKey uuid.UUID       `json:"idempotency_key"`
}

// NewOrderCreatedEvent создаёт событие из примитивных значений.
// Намеренно не принимает *order.Order, чтобы избежать циклических зависимостей.
func NewOrderCreatedEvent(
	orderID uuid.UUID,
	userID uuid.UUID,
	totalAmount decimal.Decimal,
	createdAt time.Time,
) *OrderCreatedEvent {
	return &OrderCreatedEvent{
		SchemaVersion:  "1.0",
		EventID:        uuid.New(),
		OrderID:        orderID,
		UserID:         userID,
		TotalAmount:    totalAmount,
		CreatedAt:      createdAt,
		IdempotencyKey: uuid.NewSHA1(uuid.Nil, []byte(orderID.String()+":charge")),
	}
}

// Marshal сериализует событие в JSON для публикации.
func (e *OrderCreatedEvent) Marshal() ([]byte, error) {
	res, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("events.OrderCreatedEvent.Marshal: %w", err)
	}
	return res, nil
}
