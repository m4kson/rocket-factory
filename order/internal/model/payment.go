package model

type PayOrderResponse struct {
	TransactionId string
}

type PayOrderRequest struct {
	OrderId       string
	UserId        string
	PaymentMethod PaymentMethod
}
