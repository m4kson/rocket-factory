package order_consumer

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
	"github.com/m4kson/rocket-factory/platform/pkg/kafka"
)

func (s *service) AssemblyHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.shipAssembledDecoder.Decode(msg.Value)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to decode OrderPaidRecord", slog.String("error", err.Error()))
		return err
	}

	slog.InfoContext(ctx, "Processing message",
		slog.String("topic", msg.Topic),
		slog.Any("partition", msg.Partition),
		slog.Any("offset", msg.Offset),
		slog.String("order_uuid", event.OrderID),
		slog.String("user_uuid", event.UserID),
	)

	//todo add "complete" status to order model
	orderId := uuid.MustParse(event.OrderID)
	_, err = s.orderService.UpdateStatus(ctx, orderId, model.OrderStatusUNKNOWN)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			slog.WarnContext(ctx, "order not found, skip assembled event",
				slog.String("order_uuid", event.OrderID),
				slog.Any("offset", msg.Offset),
			)
			return nil
		}

		slog.ErrorContext(ctx, "Failed to update status", slog.String("error", err.Error()))
		return err
	}

	return nil
}
