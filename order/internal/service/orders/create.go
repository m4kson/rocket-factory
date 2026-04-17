package orders

import (
	"context"

	"github.com/m4kson/rocket-factory/order/internal/model"
	repoModel "github.com/m4kson/rocket-factory/order/internal/repository/model"
)

func (s *service) CreateOrder(ctx context.Context, order model.CreateOrderRequest) (model.CreateOrderRes, error) {
	filter := model.PartsFilter{Ids: order.PartsIds}
	parts, err := s.inventoryClient.ListParts(ctx, filter)
	if err != nil {
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
