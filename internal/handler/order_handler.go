package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"e-commerce/internal/domain/order"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, dto order.CreateOrderDTO) (*order.Order, error)
	GetByID(ctx context.Context, uuid uuid.UUID) (*order.Order, error)
	ListOrders(ctx context.Context, params order.ListParams) ([]*order.Order, int64, error)
	UpdateStatus(ctx context.Context, uuid uuid.UUID, dto order.UpdateStatusDTO) (*order.Order, error)
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

// GET /orders/{id}
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

// GET /orders
func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var params order.ListParams

	if userID := query.Get("user_id"); userID != "" {
		parsedUserID, err := uuid.Parse(userID)
		if err != nil {
			h.writeError(r.Context(), w, http.StatusBadRequest, "invalid user_id format", err)
			return
		}
		params.UserID = parsedUserID
	}

	if limit := query.Get("limit"); limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err != nil {
			h.writeError(r.Context(), w, http.StatusBadRequest, "invalid limit", err)
			return
		}
		params.Limit = parsedLimit
	}

	if offset := query.Get("offset"); offset != "" {
		parsedOffset, err := strconv.Atoi(offset)
		if err != nil {
			h.writeError(r.Context(), w, http.StatusBadRequest, "invalid offset", err)
			return
		}
		params.Offset = parsedOffset
	}

	params.WithDefaults()

	orders, total, err := h.service.ListOrders(r.Context(), params)
	if err != nil {
		h.log.Error("list orders failed", "error", err)
		h.writeError(r.Context(), w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := ListOrdersResponse{
		Items:   orders,
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
		HasMore: params.Offset+params.Limit < int(total),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error("failed to encode response", "error", err)
	}
}

// PATCH /orders/id
func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		h.writeError(r.Context(), w, http.StatusBadRequest, "invalid id format", err)
		return
	}

	var dto order.UpdateStatusDTO

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.writeError(r.Context(), w, http.StatusBadRequest, "invalid json", err)
	}

	updated, err := h.service.UpdateStatus(r.Context(), orderID, dto)
	if err != nil {
		if errors.Is(err, order.ErrInvalidTransition) {
			h.writeError(r.Context(), w, http.StatusConflict, "invalid status transition", err)
			return
		}
		if errors.Is(err, order.ErrConflict) {
			h.writeError(r.Context(), w, http.StatusConflict, "resource was modified, please retry", err)
			return
		}
		if errors.Is(err, order.ErrNotFound) {
			h.writeError(r.Context(), w, http.StatusNotFound, "order not found", err)
			return
		}
		h.log.Error("update order status failed", "id", idStr, "error", err)
		h.writeError(r.Context(), w, http.StatusInternalServerError, "internal error", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(updated)
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

type ListOrdersResponse struct {
	Items   []*order.Order `json:"items"`
	Total   int64          `json:"total"`
	Limit   int            `json:"limit"`
	Offset  int            `json:"offset"`
	HasMore bool           `json:"has_more"`
}
