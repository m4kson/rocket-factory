package model

type OrderPaidEvent struct {
	ID            string
	OrderID       string
	UserID        string
	PaymentMethod string
	TransactionID string
}

type ShipAssembledEvent struct {
	ID           string
	OrderID      string
	UserID       string
	BuildTimeSec int64
}
