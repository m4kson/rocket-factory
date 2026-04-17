package service

import (
	"context"

	"github.com/m4kson/rocket-factory/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context) (*model.PayOrderResponse, error)
}
