package converter

import (
	"github.com/m4kson/rocket-factory/payment/internal/model"
	paymentV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/payment/v1"
)

func PayOrderResponseToModel(response *paymentV1.PayOrderResponse) *model.PayOrderResponse {
	return &model.PayOrderResponse{
		TransactionId: response.TransactionUuid,
	}
}
