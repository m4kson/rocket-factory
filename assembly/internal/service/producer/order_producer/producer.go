package order_producer

import (
	"context"
	"log/slog"

	"github.com/m4kson/rocket-factory/assembly/internal/model"
	"github.com/m4kson/rocket-factory/platform/pkg/kafka"
	events_v1 "github.com/m4kson/rocket-factory/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type service struct {
	orderShipAssembledProducer kafka.Producer
}

func NewService(orderShipAssembledProducer kafka.Producer) *service {
	return &service{
		orderShipAssembledProducer: orderShipAssembledProducer,
	}
}

func (p *service) ProduceOrderShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	msg := &events_v1.ShipAssembled{
		EventUuid:    event.ID,
		UserUuid:     event.UserID,
		OrderUuid:    event.OrderID,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		slog.ErrorContext(ctx, "marshal ship assemble failed: %v", slog.String("err", err.Error()))
		return err
	}

	err = p.orderShipAssembledProducer.Send(ctx, []byte(event.ID), payload)
	if err != nil {
		slog.ErrorContext(ctx, "send ship assemble failed: %v", slog.String("err", err.Error()))
		return err
	}

	return nil
}
