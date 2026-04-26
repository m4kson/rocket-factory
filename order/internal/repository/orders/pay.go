package orders

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
	repoModel "github.com/m4kson/rocket-factory/order/internal/repository/model"
	logger "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
)

func (r *repository) PayOrderById(ctx context.Context, orderId uuid.UUID, paymentMethod model.PaymentMethod, transactionId uuid.UUID) (model.PayOrderRes, error) {
	log := logger.FromContext(ctx)

	row, err := r.pool.Exec(
		ctx,
		"UPDATE orders SET payment_method = $1, status = $2, transaction_id = $3 WHERE id = $4",
		repoModel.PaymentMethod(paymentMethod),
		repoModel.OrderStatusPAID,
		transactionId,
		orderId,
	)
	if err != nil {
		log.Error("failed to pay order", slog.String("orderId", orderId.String()), slog.String("error", err.Error()))
		return model.PayOrderRes{}, fmt.Errorf("repository.PayOrderById orderId=%s: %w", orderId, err)
	}

	if row.RowsAffected() == 0 {
		log.Warn("order not found", slog.String("orderId", orderId.String()))
		return model.PayOrderRes{}, model.ErrOrderNotFound
	}

	return model.PayOrderRes{
		TransactionId: transactionId,
	}, nil
}
