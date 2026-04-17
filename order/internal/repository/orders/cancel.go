package orders

import (
	"context"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
	repoModel "github.com/m4kson/rocket-factory/order/internal/repository/model"
)

func (r *repository) CancelOrderById(ctx context.Context, orderId uuid.UUID) error {
	order, exist := r.orders[orderId.String()]
	if !exist {
		return model.ErrOrderNotFound
	}

	if order.Status == repoModel.OrderStatusPAID {

		return model.ErrOrderAlreadyPaid
	}

	order.Status = repoModel.OrderStatusCANCELLED

	r.orders[orderId.String()] = order

	return nil
}
