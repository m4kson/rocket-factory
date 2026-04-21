package orders

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
	repoModel "github.com/m4kson/rocket-factory/order/internal/repository/model"
)

func (r *repository) PayOrderById(ctx context.Context, orderId uuid.UUID, paymentMethod model.PaymentMethod, transactionId uuid.UUID) (model.PayOrderRes, error) {
	row, err := r.pool.Exec(
		ctx,
		"UPDATE orders SET payment_method = $1, status = $2, transaction_id = $3 WHERE id = $4",
		repoModel.PaymentMethod(paymentMethod),
		repoModel.OrderStatusPAID,
		transactionId,
		orderId,
	)
	if err != nil {
		return model.PayOrderRes{}, fmt.Errorf("repository.PayOrderById orderId=%s: %w", orderId, err)
	}

	if row.RowsAffected() == 0 {
		return model.PayOrderRes{}, model.ErrOrderNotFound
	}

	return model.PayOrderRes{
		TransactionId: transactionId,
	}, nil
}
