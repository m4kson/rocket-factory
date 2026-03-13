package orders

import (
	"context"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
)

func (s *service) PayOrderById(ctx context.Context, orderId uuid.UUID, paymentMethod model.PaymentMethod, userId uuid.UUID) (model.PayOrderRes, error) {
	paymentResponse, err := s.paymentClient.PayOrder(ctx, model.PayOrderRequest{
		OrderId:       orderId.String(),
		UserId:        userId.String(),
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return model.PayOrderRes{}, err
	}

	response, err := s.orderRepository.PayOrderById(ctx, orderId, paymentMethod, uuid.MustParse(paymentResponse.TransactionId))
	if err != nil {
		return model.PayOrderRes{}, err
	}

	return response, nil
}
