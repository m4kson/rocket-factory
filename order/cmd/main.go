package main

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	orderV1 "github.com/m4kson/rocket-factory/shared/pkg/openapi/order/v1"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*orderV1.GetOrderResponse //todo Тут я используя GetOrderResponse, так как OrderrDto почему-то не генерируется ogen'ом.
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[uuid.UUID]*orderV1.GetOrderResponse),
	}
}

func (s *OrderStorage) GetOrderById(orderId uuid.UUID) *orderV1.GetOrderResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderId]
	if !ok {
		return nil
	}

	return order
}

func (s *OrderStorage) CreateOrder(orderInfo *orderV1.CreateOrderRequest) *orderV1.CreateOrderResponse {
	s.mu.Lock()
	defer s.mu.Unlock()

	order := &orderV1.GetOrderResponse{
		OrderUUID:       uuid.New(),
		UserUUID:        orderInfo.UserUUID,
		PartUuids:       orderInfo.PartUuids,
		TotalPrice:      0,        //todo implement price calculation
		TransactionUUID: uuid.Nil, // todo implement transaction creation
		PaymentMethod:   orderV1.PaymentMethodUnknown,
		Status:          orderV1.OrderStatusUnknown,
	}

	s.orders[order.OrderUUID] = order

	response := orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}

	return &response
}

func (s *OrderStorage) PayOrderByUUID(orderId uuid.UUID, paymentMethod orderV1.PaymentMethod) *orderV1.PayOrderResponse {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderId]
	if !ok {
		return nil
	}

	//todo call payment service

	order.PaymentMethod = paymentMethod
	order.Status = orderV1.OrderStatusPaid

	response := orderV1.PayOrderResponse{
		TransactionID: uuid.Nil, // todo insert real transaction ID
	}

	return &response
}

func (s *OrderStorage) CancelOrderByUUID(orderId uuid.UUID) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderId]
	if !ok {
		return false
	}

	//todo call payment service to cancel transaction if order is paid

	order.Status = orderV1.OrderStatusCancelled

	return true
}

type OrderHandler struct {
	storage *OrderStorage
}

func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

//todo implement interface from oas_server_gen.go

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	return h.storage.CreateOrder(req), nil //todo implement error handling
}

func (h *OrderHandler) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	order := h.storage.GetOrderById(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "Order with UUID '" + params.OrderUUID.String() + "' not found",
		}, nil
	}

	return order, nil
}

func (h *OrderHandler) PayOrderByUUID(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderByUUIDParams) (orderV1.PayOrderByUUIDRes, error) {
	order := h.storage.PayOrderByUUID(params.OrderUUID, req.PaymentMethod)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "Order with UUID '" + params.OrderUUID.String() + "' not found",
		}, nil
	}

	return order, nil
}

func (h *OrderHandler) CancelOrderByUUID(ctx context.Context, params orderV1.CancelOrderByUUIDParams) (orderV1.CancelOrderByUUIDRes, error) {
	response := h.storage.CancelOrderByUUID(params.OrderUUID)
	if !response {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "Order with UUID '" + params.OrderUUID.String() + "' not found",
		}, nil
	}

	return &orderV1.NotFoundError{
		Code:    http.StatusNoContent,
		Message: "Order canceled successfully",
	}, nil
}
