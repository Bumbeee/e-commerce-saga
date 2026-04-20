package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr: kafka.TCP(brokers...),
			//Topic:        topic,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireOne,
			MaxAttempts:  3,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, topic string, key []byte, value []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("producer.Publish: %w", err)
	}
	return nil
}

func (p *Producer) Close() error {
	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("producer.Close: %w", err)
	}
	return nil
}
