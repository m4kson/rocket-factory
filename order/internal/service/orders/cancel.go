package orders

import (
	"context"

	"github.com/google/uuid"
)

func (s *service) CancelOrderById(ctx context.Context, orderId uuid.UUID) error {
	return s.orderRepository.CancelOrderById(ctx, orderId)
}
