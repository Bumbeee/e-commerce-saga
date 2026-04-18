package order

import "errors"

var (
	ErrNotFound          = errors.New("order not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidParams     = errors.New("invalid parameters")
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrConflict          = errors.New("resource conflict")
)
