package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
)

type OrderService interface {
	GetOrderById(ctx context.Context, orderId uuid.UUID) (model.GetOrderResponse, error)
	CreateOrder(ctx context.Context, order model.CreateOrderRequest) (model.CreateOrderRes, error)
	PayOrderById(ctx context.Context, orderId uuid.UUID, paymentMethod model.PaymentMethod, userId uuid.UUID) (model.PayOrderRes, error)
	CancelOrderById(ctx context.Context, orderId uuid.UUID) error
}
