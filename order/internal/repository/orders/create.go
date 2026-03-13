package orders

import (
	"context"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
	repoModel "github.com/m4kson/rocket-factory/order/internal/repository/model"
)

func (r *repository) CreateOrder(ctx context.Context, request repoModel.CreateOrderRequest) (model.CreateOrderRes, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	orderId := uuid.New()

	order := repoModel.Order{
		OrderId:       orderId,
		UserId:        request.UserId,
		PartsIds:      request.PartsIds,
		TotalPrice:    request.TotalPrice,
		TransactionId: request.TransactionId,
		PaymentMethod: request.PaymentMethod,
		Status:        request.Status,
	}

	r.orders[orderId.String()] = order

	return model.CreateOrderRes{
		OrderId:    orderId,
		TotalPrice: request.TotalPrice,
	}, nil
}
