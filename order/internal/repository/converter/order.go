package converter

import (
	"github.com/m4kson/rocket-factory/order/internal/model"
	repoModel "github.com/m4kson/rocket-factory/order/internal/repository/model"
)

func OrderToModel(order repoModel.Order) model.GetOrderResponse {
	return model.GetOrderResponse{
		OrderId:       order.OrderId,
		UserId:        order.UserId,
		PartsIds:      order.PartsIds,
		TotalPrice:    order.TotalPrice,
		TransactionId: order.TransactionId,
		PaymentMethod: PaymentMethodToModel(order.PaymentMethod),
		Status:        StatusToModel(order.Status),
	}
}

func PaymentMethodToModel(paymentMethod repoModel.PaymentMethod) model.PaymentMethod {
	switch paymentMethod {
	case repoModel.PaymentMethodCARD:
		return model.PaymentMethodCARD
	case repoModel.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case repoModel.PaymentMethodCREDITCARD:
		return model.PaymentMethodCREDITCARD
	case repoModel.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodINVESTORMONEY
	default:
		return model.PaymentMethodUNKNOWN
	}
}

func StatusToModel(status repoModel.OrderStatus) model.OrderStatus {
	switch status {
	case repoModel.OrderStatusPENDINGPAYMENT:
		return model.OrderStatusPENDINGPAYMENT
	case repoModel.OrderStatusPAID:
		return model.OrderStatusPAID
	case repoModel.OrderStatusCANCELLED:
		return model.OrderStatusCANCELLED
	default:
		return model.OrderStatusUNKNOWN
	}
}
