package producer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
)

type producer struct {
	syncProducer sarama.SyncProducer
	topic        string
	logger       *slog.Logger
}

func NewProducer(syncProducer sarama.SyncProducer, topic string, logger *slog.Logger) *producer {
	return &producer{
		syncProducer: syncProducer,
		topic:        topic,
		logger:       logger,
	}
}

func (p *producer) Send(ctx context.Context, key, value []byte) error {
	partition, offset, err := p.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	})
	if err != nil {
		p.logger.Error("Failed to send message", slog.String("topic", p.topic), slog.String("key", string(key)), slog.String("error", err.Error()))
		return fmt.Errorf("failed to send message: %w", err)
	}

	p.logger.Info("Message sent",
		slog.String("topic", p.topic),
		slog.Int("partition", int(partition)),
		slog.Int64("offset", offset),
		slog.String("key", string(key)),
		slog.String("value", string(value)),
	)

	return nil
}
