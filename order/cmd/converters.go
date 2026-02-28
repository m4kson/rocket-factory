package main

import (
	orderV1 "github.com/m4kson/rocket-factory/shared/pkg/openapi/order/v1"
	paymentV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/payment/v1"
)

// ConvertOpenAPIPaymentMethodToProto преобразует OpenAPI PaymentMethod в proto PaymentMethod
// OpenAPI PaymentMethod - это строка (UNKNOWN, CARD, SBP, CREDIT_CARD, INVESTOR_MONEY)
// Proto PaymentMethod - это oneof с вложенными типами
func ConvertOpenAPIPaymentMethodToProto(method orderV1.PaymentMethod) *paymentV1.PaymentMethod {
	protoMethod := &paymentV1.PaymentMethod{}

	switch method {
	case orderV1.PaymentMethodCARD:
		protoMethod.Method = &paymentV1.PaymentMethod_Card{Card: "card"}
	case orderV1.PaymentMethodSBP:
		protoMethod.Method = &paymentV1.PaymentMethod_Sbp{Sbp: "sbp"}
	case orderV1.PaymentMethodCREDITCARD:
		protoMethod.Method = &paymentV1.PaymentMethod_CreditCard{CreditCard: "credit_card"}
	case orderV1.PaymentMethodINVESTORMONEY:
		protoMethod.Method = &paymentV1.PaymentMethod_InvestorMoney{InvestorMoney: "investor_money"}
	case orderV1.PaymentMethodUNKNOWN:
		fallthrough
	default:
		protoMethod.Method = &paymentV1.PaymentMethod_Unknown{Unknown: "unknown"}
	}

	return protoMethod
}

// ConvertProtoPaymentMethodToOpenAPI преобразует proto PaymentMethod в OpenAPI PaymentMethod
func ConvertProtoPaymentMethodToOpenAPI(method *paymentV1.PaymentMethod) orderV1.PaymentMethod {
	if method == nil {
		return orderV1.PaymentMethodUNKNOWN
	}

	switch method.Method.(type) {
	case *paymentV1.PaymentMethod_Card:
		return orderV1.PaymentMethodCARD
	case *paymentV1.PaymentMethod_Sbp:
		return orderV1.PaymentMethodSBP
	case *paymentV1.PaymentMethod_CreditCard:
		return orderV1.PaymentMethodCREDITCARD
	case *paymentV1.PaymentMethod_InvestorMoney:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}
