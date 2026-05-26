package order_producer

import (
	"context"
	"log/slog"

	"github.com/m4kson/rocket-factory/order/internal/model"
	"github.com/m4kson/rocket-factory/platform/pkg/kafka"
	events_v1 "github.com/m4kson/rocket-factory/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type service struct {
	orderPaidProducer kafka.Producer
}

func NewService(orderPaidProducer kafka.Producer) *service {
	return &service{orderPaidProducer: orderPaidProducer}
}

func (p *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	msg := &events_v1.OrderPaid{
		EventUuid:       event.ID,
		UserUuid:        event.UserID,
		OrderUuid:       event.OrderID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUuid: event.TransactionID,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		slog.ErrorContext(ctx, "marshal order paid failed: %v", slog.String("err", err.Error()))
		return err
	}

	err = p.orderPaidProducer.Send(ctx, []byte(event.ID), payload)
	if err != nil {
		slog.ErrorContext(ctx, "send order paid failed: %v", slog.String("err", err.Error()))
		return err
	}

	return nil
}
