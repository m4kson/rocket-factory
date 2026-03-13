package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/m4kson/rocket-factory/order/internal/model"
	orderV1 "github.com/m4kson/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrderById(ctx context.Context, params orderV1.CancelOrderByUUIDParams) (orderV1.CancelOrderByUUIDRes, error) {
	err := a.orderService.CancelOrderById(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Order with UUID '" + params.OrderUUID.String() + "' not found",
			}, nil
		}
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusConflict,
				Message: "Order with UUID '" + params.OrderUUID.String() + "' already paid and cannot be cancelled",
			}, nil
		}

		return nil, err
	}

	return &orderV1.CancelOrderByUUIDNoContent{}, nil
}
