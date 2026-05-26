package order_consumer

import (
	"context"
	"log/slog"

	kafkaConverter "github.com/m4kson/rocket-factory/order/internal/converter/kafka"
	def "github.com/m4kson/rocket-factory/order/internal/service"
	"github.com/m4kson/rocket-factory/platform/pkg/kafka"
)

var _ def.OrderConsumerService = (*service)(nil)

type service struct {
	shipAssembledConsumer kafka.Consumer
	shipAssembledDecoder  kafkaConverter.OrderAssembledDecoder
	orderService          def.OrderService
}

func NewService(shipAssembledConsumer kafka.Consumer, shipAssembledDecoder kafkaConverter.OrderAssembledDecoder, orderService def.OrderService) *service {
	return &service{
		shipAssembledConsumer: shipAssembledConsumer,
		shipAssembledDecoder:  shipAssembledDecoder,
		orderService:          orderService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	slog.InfoContext(ctx, "Starting order consumer service")

	err := s.shipAssembledConsumer.Consume(ctx, s.AssemblyHandler)
	if err != nil {
		slog.ErrorContext(ctx, "Consume from ufo.recorded topic error", slog.String("error", err.Error()))
		return err
	}

	return nil
}
