package order

import (
	"github.com/google/uuid"
)

type ListParams struct {
	Limit  int
	Offset int
	UserID uuid.UUID
}

func (p *ListParams) WithDefaults() {
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}
