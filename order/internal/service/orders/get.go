package orders

import (
	"context"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
)

func (s *service) GetOrderById(ctx context.Context, orderId uuid.UUID) (model.GetOrderResponse, error) {
	order, err := s.orderRepository.GetOrderById(ctx, orderId)
	if err != nil {
		return model.GetOrderResponse{}, err
	}

	return order, nil
}
