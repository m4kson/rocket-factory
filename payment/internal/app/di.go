package app

import (
	paymentV1API "github.com/m4kson/rocket-factory/payment/internal/api/payment/v1"
	service "github.com/m4kson/rocket-factory/payment/internal/service"
	paymentService "github.com/m4kson/rocket-factory/payment/internal/service/payment"
	paymentV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API paymentV1.PaymentServiceServer

	paymentService service.PaymentService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API() paymentV1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = paymentV1API.NewAPI(d.PaymentService())
	}

	return d.paymentV1API
}

func (d *diContainer) PaymentService() service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewPaymentService()
	}

	return d.paymentService
}
