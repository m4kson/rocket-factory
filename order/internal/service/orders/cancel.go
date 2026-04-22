package orders

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
)

func (s *service) CancelOrderById(ctx context.Context, orderId uuid.UUID) error {
	order, err := s.orderRepository.GetOrderById(ctx, orderId)
	if err != nil {
		return fmt.Errorf("service.CancelOrderById orderId=%s: %w", orderId, err)
	}

	if order.Status == model.OrderStatusPAID {
		return model.ErrOrderAlreadyPaid
	}

	return s.orderRepository.CancelOrderById(ctx, orderId)
}
