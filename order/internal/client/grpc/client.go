package grpc

import (
	"context"

	"github.com/m4kson/rocket-factory/order/internal/model"
)

type PaymentClient interface {
	PayOrder(ctx context.Context, requestModel model.PayOrderRequest) (model.PayOrderResponse, error)
}

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
