package orders

import (
	"context"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
	repoModel "github.com/m4kson/rocket-factory/order/internal/repository/model"
)

func (r *repository) PayOrderById(ctx context.Context, orderId uuid.UUID, paymentMethod model.PaymentMethod, transactionId uuid.UUID) (model.PayOrderRes, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order := r.orders[orderId.String()]

	order.PaymentMethod = repoModel.PaymentMethod(paymentMethod)
	order.Status = repoModel.OrderStatusPAID
	order.TransactionId = &transactionId

	r.orders[orderId.String()] = order

	return model.PayOrderRes{
		TransactionId: transactionId,
	}, nil
}
