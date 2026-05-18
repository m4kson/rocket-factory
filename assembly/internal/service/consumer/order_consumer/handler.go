package order_consumer

import (
	"context"
	"log/slog"
	"time"

	"github.com/m4kson/rocket-factory/assembly/internal/model"
	"github.com/m4kson/rocket-factory/platform/pkg/kafka"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
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
		slog.String("transaction_uuid", event.TransactionID))

	time.Sleep(10 * time.Second)

	shipAssembled := model.ShipAssembledEvent{
		ID:           event.ID,
		OrderID:      event.OrderID,
		UserID:       event.UserID,
		BuildTimeSec: int64(10),
	}

	err = s.orderShipAssembledProducer.ProduceOrderShipAssembled(ctx, shipAssembled)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to produce ship assembled event",
			slog.String("event_uuid", shipAssembled.ID),
			slog.String("order_uuid", shipAssembled.OrderID),
			slog.Int64("build_time_sec", shipAssembled.BuildTimeSec),
			slog.String("error", err.Error()),
		)
		return err
	}

	return nil
}
