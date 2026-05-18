package kafka

import (
	"context"
	"log/slog"

	"github.com/m4kson/rocket-factory/platform/pkg/kafka"
	"github.com/m4kson/rocket-factory/platform/pkg/kafka/consumer"
)

func Logging(logger *slog.Logger) consumer.Middleware {
	return func(next kafka.MessageHandler) kafka.MessageHandler {
		return func(ctx context.Context, msg kafka.Message) error {
			logger.InfoContext(ctx, "Kafka msg received", slog.String("topic", msg.Topic))
			return next(ctx, msg)
		}
	}
}
