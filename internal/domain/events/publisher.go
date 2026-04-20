package events

import (
	"context"
)

type Publisher interface {
	Publish(ctx context.Context, topic string, key []byte, value []byte) error
	Close() error
}
