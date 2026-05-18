package order_consumer

import (
	"context"
	"log/slog"

	kafkaConverter "github.com/m4kson/rocket-factory/assembly/internal/converter/kafka"
	assemblyService "github.com/m4kson/rocket-factory/assembly/internal/service"
	def "github.com/m4kson/rocket-factory/assembly/internal/service"
	"github.com/m4kson/rocket-factory/platform/pkg/kafka"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderPaidConsumer          kafka.Consumer
	orderPaidDecoder           kafkaConverter.OrderPaidDecoder
	orderShipAssembledProducer assemblyService.ShipAssembledProducerService
}

func NewService(orderPaidConsumer kafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder, orderShipAssembledProducer assemblyService.ShipAssembledProducerService) *service {
	return &service{
		orderPaidConsumer:          orderPaidConsumer,
		orderPaidDecoder:           orderPaidDecoder,
		orderShipAssembledProducer: orderShipAssembledProducer,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	slog.InfoContext(ctx, "Starting orderPaid consumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to consume order", slog.String("err", err.Error()))
		return err
	}

	return nil
}
