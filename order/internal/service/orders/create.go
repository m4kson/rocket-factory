package orders

import (
	"context"
	"log/slog"

	"github.com/m4kson/rocket-factory/order/internal/model"
	repoModel "github.com/m4kson/rocket-factory/order/internal/repository/model"
	logger "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
)

func (s *service) CreateOrder(ctx context.Context, order model.CreateOrderRequest) (model.CreateOrderRes, error) {
	log := logger.FromContext(ctx)

	filter := model.PartsFilter{Ids: order.PartsIds}
	parts, err := s.inventoryClient.ListParts(ctx, filter)
	if err != nil {
		log.Error("failed to list parts", slog.String("error", err.Error()))
		return model.CreateOrderRes{}, err
	}

	var totalPrice float32
	for _, part := range parts {
		totalPrice += part.Price
	}

	request := repoModel.CreateOrderRequest{
		UserId:        order.UserId,
		PartsIds:      order.PartsIds,
		TotalPrice:    totalPrice,
		TransactionId: nil,
		PaymentMethod: repoModel.PaymentMethodUNKNOWN,
		Status:        repoModel.OrderStatusPENDINGPAYMENT,
	}

	response, err := s.orderRepository.CreateOrder(ctx, request)
	if err != nil {
		return model.CreateOrderRes{}, err
	}

	return response, nil
}
