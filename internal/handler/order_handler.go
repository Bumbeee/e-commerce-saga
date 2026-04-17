package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"e-commerce/internal/domain/order"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, dto order.CreateOrderDTO) (*order.Order, error)
	GetByID(ctx context.Context, uuid uuid.UUID) (*order.Order, error)
}

type OrderHandler struct {
	service OrderService
	log     *slog.Logger
}

func NewOrderHandler(service OrderService, log *slog.Logger) *OrderHandler {
	return &OrderHandler{service: service, log: log}
}

// POST /orders
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto order.CreateOrderDTO

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.writeError(r.Context(), w, http.StatusBadRequest, "invalid json", err)
		return
	}

	ord, err := h.service.CreateOrder(r.Context(), dto)
	if err != nil {
		if errors.Is(err, order.ErrInvalidInput) {
			h.writeError(r.Context(), w, http.StatusBadRequest, "validation failed", err)
			return
		}
		h.log.Error("create order failed", "error", err)
		h.writeError(r.Context(), w, http.StatusInternalServerError, "internal error", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/orders/"+ord.ID.String())
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(ord); err != nil {
		h.log.Error("failed to encode response", "error", err)
	}
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		h.writeError(r.Context(), w, http.StatusBadRequest, "invalid id format", err)
		return
	}

	ord, err := h.service.GetByID(r.Context(), orderID)
	if err != nil {
		if errors.Is(err, order.ErrNotFound) {
			h.writeError(r.Context(), w, http.StatusNotFound, "order not found", err)
			return
		}
		h.log.Error("get order failed", "id", idStr, "error", err)
		h.writeError(r.Context(), w, http.StatusInternalServerError, "internal server error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(ord); err != nil {
		h.log.Error("failed to encode response", "error", err)
	}

}

func (h *OrderHandler) writeError(ctx context.Context, w http.ResponseWriter, status int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := ErrorResponse{Error: message}
	if err != nil && h.log.Enabled(ctx, slog.LevelDebug) {
		response.Details = err.Error()
	}

	_ = json.NewEncoder(w).Encode(response)
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
