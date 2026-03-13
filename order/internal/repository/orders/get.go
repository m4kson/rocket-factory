package orders

import (
	"context"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
	"github.com/m4kson/rocket-factory/order/internal/repository/converter"
)

func (r *repository) GetOrderById(ctx context.Context, orderId uuid.UUID) (model.GetOrderResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exist := r.orders[orderId.String()]
	if !exist {
		return model.GetOrderResponse{}, model.ErrOrderNotFound
	}

	response := converter.OrderToModel(order)

	return response, nil
}
