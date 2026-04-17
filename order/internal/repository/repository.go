package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
	repoModel "github.com/m4kson/rocket-factory/order/internal/repository/model"
)

type OrderRepository interface {
	GetOrderById(ctx context.Context, orderId uuid.UUID) (model.GetOrderResponse, error)
	CreateOrder(ctx context.Context, order repoModel.CreateOrderRequest) (model.CreateOrderRes, error)
	PayOrderById(ctx context.Context, orderId uuid.UUID, paymentMethod model.PaymentMethod, transactionId uuid.UUID) (model.PayOrderRes, error)
	CancelOrderById(ctx context.Context, orderId uuid.UUID) error
}
