package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/m4kson/rocket-factory/order/internal/converter"
	"github.com/m4kson/rocket-factory/order/internal/model"
	orderV1 "github.com/m4kson/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderById(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	order, err := a.orderService.GetOrderById(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Order with UUID '" + params.OrderUUID.String() + "' not found",
			}, nil
		}
		return nil, err
	}

	response := converter.GetOrderResponseToDto(order)
	return response, nil
}
