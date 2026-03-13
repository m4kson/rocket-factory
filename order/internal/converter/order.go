package converter

import (
	"github.com/m4kson/rocket-factory/order/internal/model"
	orderV1 "github.com/m4kson/rocket-factory/shared/pkg/openapi/order/v1"
)

func CreateOrderRequestToModel(req *orderV1.CreateOrderRequest) model.CreateOrderRequest {
	return model.CreateOrderRequest{
		UserId:   req.UserUUID,
		PartsIds: req.PartUuids,
	}
}

func PaymentMethodToModel(paymentMethod orderV1.PaymentMethod) model.PaymentMethod {
	switch paymentMethod {
	case orderV1.PaymentMethodCARD:
		return model.PaymentMethodCARD
	case orderV1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case orderV1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCREDITCARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodINVESTORMONEY
	default:
		return model.PaymentMethodUNKNOWN
	}
}

func GetOrderResponseToDto(model model.GetOrderResponse) *orderV1.GetOrderResponse {
	var transactionId orderV1.OptUUID
	if model.TransactionId != nil {
		transactionId = orderV1.NewOptUUID(*model.TransactionId)
	}

	return &orderV1.GetOrderResponse{
		OrderUUID:       model.OrderId,
		UserUUID:        model.UserId,
		PartUuids:       model.PartsIds,
		TotalPrice:      model.TotalPrice,
		TransactionUUID: transactionId,
		PaymentMethod:   PaymentMethodToDto(model.PaymentMethod),
	}
}

func PaymentMethodToDto(method model.PaymentMethod) orderV1.PaymentMethod {
	switch method {
	case model.PaymentMethodUNKNOWN:
		return orderV1.PaymentMethodUNKNOWN
	case model.PaymentMethodCARD:
		return orderV1.PaymentMethodCARD
	case model.PaymentMethodCREDITCARD:
		return orderV1.PaymentMethodCREDITCARD
	case model.PaymentMethodSBP:
		return orderV1.PaymentMethodSBP
	case model.PaymentMethodINVESTORMONEY:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}
