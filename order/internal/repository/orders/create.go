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

func (r *repository) CreateOrder(ctx context.Context, request repoModel.CreateOrderRequest) (model.CreateOrderRes, error) {
	log := logger.FromContext(ctx)

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

	_, err := r.pool.Exec(
		ctx,
		"INSERT INTO orders (id, user_id, parts_ids, total_price, transaction_id, payment_method, status) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		order.OrderId,
		order.UserId,
		order.PartsIds,
		order.TotalPrice,
		order.TransactionId,
		order.PaymentMethod,
		order.Status,
	)
	if err != nil {
		log.Error("failed to create order", slog.String("error", err.Error()))
		return model.CreateOrderRes{}, fmt.Errorf("repository.CreateOrder: %w", err)
	}

	return model.CreateOrderRes{
		OrderId:    orderId,
		TotalPrice: request.TotalPrice,
	}, nil
}
