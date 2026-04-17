package model

import "github.com/google/uuid"

type Order struct {
	OrderId       uuid.UUID
	UserId        uuid.UUID
	PartsIds      []uuid.UUID
	TotalPrice    float32
	TransactionId *uuid.UUID
	PaymentMethod PaymentMethod
	Status        OrderStatus
}

type PaymentMethod string

const (
	PaymentMethodUNKNOWN       PaymentMethod = "UNKNOWN"
	PaymentMethodCARD          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCREDITCARD    PaymentMethod = "CREDIT_CARD"
	PaymentMethodINVESTORMONEY PaymentMethod = "INVESTOR_MONEY"
)

type OrderStatus string

const (
	OrderStatusUNKNOWN        OrderStatus = "UNKNOWN"
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
)

type CreateOrderRequest struct {
	UserId        uuid.UUID
	PartsIds      []uuid.UUID
	TotalPrice    float32
	TransactionId *uuid.UUID
	PaymentMethod PaymentMethod
	Status        OrderStatus
}
