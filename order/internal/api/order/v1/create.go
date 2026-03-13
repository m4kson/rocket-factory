package v1

import (
	"context"

	"github.com/m4kson/rocket-factory/order/internal/converter"
	orderV1 "github.com/m4kson/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	request := converter.CreateOrderRequestToModel(req)
	order, err := a.orderService.CreateOrder(ctx, request)
	if err != nil {
		return nil, err //todo добавить проверку всех PartsId на существование, если хотя бы одной нет - вернуть ошибку
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderId,
		TotalPrice: order.TotalPrice,
	}, nil
}
