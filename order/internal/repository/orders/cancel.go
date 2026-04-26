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

func (r *repository) CancelOrderById(ctx context.Context, orderId uuid.UUID) error {
	log := logger.FromContext(ctx)

	row, err := r.pool.Exec(ctx,
		"UPDATE orders SET status = $1 WHERE id = $2 AND status != $3",
		repoModel.OrderStatusCANCELLED,
		orderId,
		repoModel.OrderStatusPAID,
	)
	if err != nil {
		log.Error("failed to cancel order", slog.String("orderId", orderId.String()), slog.String("error", err.Error()))
		return fmt.Errorf("repository.CancelOrderById orderId=%s: %w", orderId, err)
	}

	if row.RowsAffected() == 0 {
		log.Warn("order not found", slog.String("orderId", orderId.String()))
		return model.ErrOrderNotFound
	}

	return nil
}
