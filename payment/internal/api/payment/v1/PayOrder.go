package v1

import (
	"context"

	paymentV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, _ *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	response, err := a.paymentService.PayOrder(ctx)
	if err != nil {
		return nil, err
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: response.TransactionId,
	}, nil
}
