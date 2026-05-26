package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/converter"
	"github.com/m4kson/rocket-factory/order/internal/model"
	orderV1 "github.com/m4kson/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrderByUUID(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderByUUIDParams) (orderV1.PayOrderByUUIDRes, error) {
	//userId := ctx.Value("user_id").(string)
	//todo put userId in context, next string is zaglushka
	userId := "123e4567-e89b-12d3-a456-426614174000"
	transaction, err := a.orderService.PayOrderById(ctx, params.OrderUUID, converter.PaymentMethodToModel(req.PaymentMethod), uuid.MustParse(userId))
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Order with UUID '" + params.OrderUUID.String() + "' not found",
			}, nil
		}

		return nil, err
	}
	return &orderV1.PayOrderResponse{
		TransactionID: transaction.TransactionId,
	}, nil
}
