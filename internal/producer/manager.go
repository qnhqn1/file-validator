package producer

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/qnhqn1/file-validator/config"
)


type Manager struct {
	writer *kafka.Writer
	cfg    *config.Config
}


func New(cfg *config.Config) *Manager {
	w := &kafka.Writer{
		Addr:  kafka.TCP(cfg.Kafka.Brokers...),
		Topic: cfg.Topics.Response,
		Async: false,
	}
	return &Manager{writer: w, cfg: cfg}
}


func (m *Manager) SendValidated(ctx context.Context, key []byte, value []byte) error {
	msg := kafka.Message{Key: key, Value: value}
	if err := m.writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("производитель: ошибка записи: %v", err)
		return err
	}
	return nil
}


func (m *Manager) Close() {
	if m.writer != nil {
		_ = m.writer.Close()
	}
}


