package kafka

import "context"

// NoOpPublisher — заглушка, реализующая интерфейс events.Publisher.
// Используется, когда Kafka не настроена (локальная разработка / тесты).
// Позволяет избежать проверок на nil в бизнес-логике.
type NoOpPublisher struct{}

// Publish имитирует успешную доставку события.
// Для NoOpPublisher всегда возвращает nil (успех).
func (n *NoOpPublisher) Publish(ctx context.Context, topic string, key []byte, value []byte) error {
	return nil
}

// Close имитирует корректное освобождение ресурсов.
func (n *NoOpPublisher) Close() error {
	return nil
}
