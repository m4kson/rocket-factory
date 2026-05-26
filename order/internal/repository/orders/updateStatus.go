package orders

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/m4kson/rocket-factory/order/internal/model"
	"github.com/m4kson/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/m4kson/rocket-factory/order/internal/repository/model"
	logger "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
)

func (r *repository) UpdateStatus(ctx context.Context, orderId uuid.UUID, status model.OrderStatus) (*model.GetOrderResponse, error) {
	log := logger.FromContext(ctx)

	var row repoModel.Order

	err := r.pool.QueryRow(
		ctx,
		`UPDATE orders SET status = $1 WHERE id = $2
		 RETURNING id, user_id, parts_ids, total_price, transaction_id, payment_method, status`,
		repoModel.OrderStatus(status),
		orderId,
	).Scan(
		&row.OrderId,
		&row.UserId,
		&row.PartsIds,
		&row.TotalPrice,
		&row.TransactionId,
		&row.PaymentMethod,
		&row.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.WarnContext(ctx, "order not found", slog.String("orderId", orderId.String()))
			return nil, model.ErrOrderNotFound
		}

		log.ErrorContext(ctx, "failed to update order status",
			slog.String("order_id", orderId.String()),
			slog.String("status", string(status)),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	resp := converter.OrderToModel(row)
	return &resp, nil
}
