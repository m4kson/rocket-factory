package orders

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
	logger "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
)

func (s *service) PayOrderById(ctx context.Context, orderId uuid.UUID, paymentMethod model.PaymentMethod, userId uuid.UUID) (model.PayOrderRes, error) {
	log := logger.FromContext(ctx)

	paymentResponse, err := s.paymentClient.PayOrder(ctx, model.PayOrderRequest{
		OrderId:       orderId.String(),
		UserId:        userId.String(),
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		log.Error("failed to pay order", slog.String("orderId", orderId.String()), slog.String("error", err.Error()))
		return model.PayOrderRes{}, err
	}

	response, err := s.orderRepository.PayOrderById(ctx, orderId, paymentMethod, uuid.MustParse(paymentResponse.TransactionId))
	if err != nil {
		return model.PayOrderRes{}, err
	}

	return response, nil
}
