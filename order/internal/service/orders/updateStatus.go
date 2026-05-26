package orders

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
	logger "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
)

func (s *service) UpdateStatus(ctx context.Context, orderId uuid.UUID, status model.OrderStatus) (*model.GetOrderResponse, error) {
	log := logger.FromContext(ctx)

	order, err := s.orderRepository.GetOrderById(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by id %w", err)
	}

	if order.Status == status {
		log.ErrorContext(ctx, "order already has this status", slog.String("order_id", orderId.String()), slog.String("status", string(status)))
		return &order, nil
	}

	response, err := s.orderRepository.UpdateStatus(ctx, orderId, status)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status %w", err)
	}

	return response, nil
}
