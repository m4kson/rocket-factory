package consumer

import (
	"context"
	"errors"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/m4kson/rocket-factory/platform/pkg/kafka"
)

type consumer struct {
	group       sarama.ConsumerGroup
	topics      []string
	logger      *slog.Logger
	middlewares []Middleware
}

func NewConsumer(group sarama.ConsumerGroup, topics []string, logger *slog.Logger, middlewares ...Middleware) *consumer {
	return &consumer{
		group:       group,
		topics:      topics,
		logger:      logger,
		middlewares: middlewares,
	}
}

func (c *consumer) Consume(ctx context.Context, handler kafka.MessageHandler) error {
	newGroupHandler := NewGroupHandler(handler, c.logger, c.middlewares...)

	for {
		if err := c.group.Consume(ctx, c.topics, newGroupHandler); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			c.logger.Error("Kafka consume error", slog.String("error", err.Error()))
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		c.logger.Info("Kafka consumer group rebalancing...")
	}
}
