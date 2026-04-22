package orders

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
	repoModel "github.com/m4kson/rocket-factory/order/internal/repository/model"
)

func (r *repository) CancelOrderById(ctx context.Context, orderId uuid.UUID) error {
	row, err := r.pool.Exec(ctx,
		"UPDATE orders SET status = $1 WHERE id = $2 AND status != $3",
		repoModel.OrderStatusCANCELLED,
		orderId,
		repoModel.OrderStatusPAID,
	)
	if err != nil {
		return fmt.Errorf("repository.CancelOrderById orderId=%s: %w", orderId, err)
	}

	if row.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}
