package converter

import (
	"github.com/m4kson/rocket-factory/order/internal/model"
	paymentV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/payment/v1"
)

func PayOrderRequestToProto(model model.PayOrderRequest) paymentV1.PayOrderRequest {
	return paymentV1.PayOrderRequest{
		OrderUuid:     model.OrderId,
		UserUuid:      model.UserId,
		PaymentMethod: PaymentMethodToProto(model.PaymentMethod),
	}
}

func PayOrderResponseToModel(proto *paymentV1.PayOrderResponse) model.PayOrderResponse {
	return model.PayOrderResponse{
		TransactionId: proto.TransactionUuid,
	}
}

func PaymentMethodToProto(method model.PaymentMethod) *paymentV1.PaymentMethod {
	protoMethod := &paymentV1.PaymentMethod{}

	switch method {
	case model.PaymentMethodCARD:
		protoMethod.Method = &paymentV1.PaymentMethod_Card{Card: "card"}
	case model.PaymentMethodSBP:
		protoMethod.Method = &paymentV1.PaymentMethod_Sbp{Sbp: "sbp"}
	case model.PaymentMethodCREDITCARD:
		protoMethod.Method = &paymentV1.PaymentMethod_CreditCard{CreditCard: "credit_card"}
	case model.PaymentMethodINVESTORMONEY:
		protoMethod.Method = &paymentV1.PaymentMethod_InvestorMoney{InvestorMoney: "investor_money"}
	case model.PaymentMethodUNKNOWN:
		fallthrough
	default:
		protoMethod.Method = &paymentV1.PaymentMethod_Unknown{Unknown: "unknown"}
	}

	return protoMethod
}
