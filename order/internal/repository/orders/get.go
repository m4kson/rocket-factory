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

func (r *repository) GetOrderById(ctx context.Context, orderId uuid.UUID) (model.GetOrderResponse, error) {
	log := logger.FromContext(ctx)

	var row repoModel.Order

	err := r.pool.QueryRow(ctx,
		"SELECT id, user_id, parts_ids, total_price, transaction_id, payment_method, status FROM orders WHERE id = $1",
		orderId).Scan(
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
			log.Warn("order not found", slog.String("orderId", orderId.String()))
			return model.GetOrderResponse{}, model.ErrOrderNotFound
		}

		log.Error("failed to get order by id", slog.String("orderId", orderId.String()), slog.String("error", err.Error()))
		return model.GetOrderResponse{}, fmt.Errorf("repository.GetOrderById orderId=%s: %w", orderId, err)
	}

	return converter.OrderToModel(row), nil
}
