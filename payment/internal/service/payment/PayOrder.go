package payment

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/payment/internal/model"
)

func (s *service) PayOrder(_ context.Context) (*model.PayOrderResponse, error) {
	transactionId := uuid.New().String()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionId)

	return &model.PayOrderResponse{
		TransactionId: transactionId,
	}, nil
}
