package orders

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
	logger "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
)

func (s *service) CancelOrderById(ctx context.Context, orderId uuid.UUID) error {
	log := logger.FromContext(ctx)

	order, err := s.orderRepository.GetOrderById(ctx, orderId)
	if err != nil {
		return fmt.Errorf("service.CancelOrderById orderId=%s: %w", orderId, err)
	}

	if order.Status == model.OrderStatusPAID {
		log.Warn("order is paid", slog.String("orderId", orderId.String()))
		return model.ErrOrderAlreadyPaid
	}

	return s.orderRepository.CancelOrderById(ctx, orderId)
}
