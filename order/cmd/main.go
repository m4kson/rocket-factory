package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/m4kson/rocket-factory/order/internal/client/grpc/inventory"
	"github.com/m4kson/rocket-factory/order/internal/client/grpc/payment"
	ordersService "github.com/m4kson/rocket-factory/order/internal/service/orders"
	ordersRepo "github.com/m4kson/rocket-factory/order/internal/repository/orders"
	ordersApi "github.com/m4kson/rocket-factory/order/internal/api/order/v1"
	orderV1 "github.com/m4kson/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

//type OrderStorage struct {
//	mu     sync.RWMutex
//	orders map[uuid.UUID]*OrderDTO
//
//	paymentClient   *PaymentClient
//	inventoryClient *InventoryClient
//}
//
//type OrderDTO struct {
//	OrderUUID       uuid.UUID
//	UserUUID        uuid.UUID
//	PartUuids       []uuid.UUID
//	TotalPrice      float32
//	TransactionUUID uuid.UUID
//	PaymentMethod   orderV1.PaymentMethod
//	Status          orderV1.OrderStatus
//}
//
//func NewOrderStorage(paymentClient *PaymentClient, inventoryClient *InventoryClient) *OrderStorage {
//	return &OrderStorage{
//		orders:          make(map[uuid.UUID]*OrderDTO),
//		paymentClient:   paymentClient,
//		inventoryClient: inventoryClient,
//	}
//}
//
//func buildOrderDTOFromCreate(req *orderV1.CreateOrderRequest, totalPrice float32) *OrderDTO {
//	return &OrderDTO{
//		OrderUUID:       uuid.New(),
//		UserUUID:        req.UserUUID,
//		PartUuids:       req.PartUuids,
//		TotalPrice:      totalPrice,
//		TransactionUUID: uuid.Nil,
//		PaymentMethod:   orderV1.PaymentMethodUNKNOWN,
//		Status:          orderV1.OrderStatusPENDINGPAYMENT,
//	}
//}
//
//func buildGetOrderResponseFromDTO(d *OrderDTO) *orderV1.GetOrderResponse {
//	return &orderV1.GetOrderResponse{
//		OrderUUID:       d.OrderUUID,
//		UserUUID:        d.UserUUID,
//		PartUuids:       d.PartUuids,
//		TotalPrice:      d.TotalPrice,
//		TransactionUUID: d.TransactionUUID,
//		PaymentMethod:   d.PaymentMethod,
//		Status:          d.Status,
//	}
//}
//
//func (s *OrderStorage) GetOrderById(orderId uuid.UUID) *orderV1.GetOrderResponse {
//	s.mu.RLock()
//	defer s.mu.RUnlock()
//
//	order, ok := s.orders[orderId]
//	if !ok {
//		return nil
//	}
//
//	return buildGetOrderResponseFromDTO(order)
//}
//
//func (s *OrderStorage) CreateOrder(orderInfo *orderV1.CreateOrderRequest) *orderV1.CreateOrderResponse {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	totalPrice := float32(0.0)
//	for _, partId := range orderInfo.PartUuids {
//		part, err := GetPart(context.Background(), s.inventoryClient.client, &inventoryV1.GetPartRequest{PartUuid: partId.String()})
//		if err != nil {
//			log.Printf("Failed to get part %s: %v", partId, err)
//			continue
//		}
//		totalPrice += part.Part.Price
//	}
//
//	dto := buildOrderDTOFromCreate(orderInfo, totalPrice)
//	s.orders[dto.OrderUUID] = dto
//
//	response := orderV1.CreateOrderResponse{
//		OrderUUID:  dto.OrderUUID,
//		TotalPrice: dto.TotalPrice,
//	}
//
//	return &response
//}
//
//func (s *OrderStorage) PayOrderByUUID(ctx context.Context, orderId uuid.UUID, paymentMethod orderV1.PaymentMethod) *orderV1.PayOrderResponse {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	dto, ok := s.orders[orderId]
//	if !ok {
//		return nil
//	}
//
//	paymentRequest := &paymentV1.PayOrderRequest{
//		OrderUuid:     orderId.String(),
//		UserUuid:      dto.UserUUID.String(),
//		PaymentMethod: ConvertOpenAPIPaymentMethodToProto(paymentMethod),
//	}
//	transactionUUID, err := PayOrder(ctx, s.paymentClient.client, paymentRequest)
//	if err != nil {
//		log.Printf("Failed to process payment for order %s: %v", orderId, err)
//		return nil
//	}
//
//	dto.PaymentMethod = paymentMethod
//	dto.Status = orderV1.OrderStatusPAID
//	dto.TransactionUUID = uuid.MustParse(transactionUUID.TransactionUuid)
//
//	response := orderV1.PayOrderResponse{
//		TransactionID: dto.TransactionUUID,
//	}
//
//	return &response
//}
//
//func (s *OrderStorage) CancelOrderByUUID(orderId uuid.UUID) bool {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	dto, ok := s.orders[orderId]
//	if !ok {
//		return false
//	}
//
//	dto.Status = orderV1.OrderStatusCANCELLED
//
//	return true
//}
//
//type OrderHandler struct {
//	storage *OrderStorage
//}
//
//func NewOrderHandler(storage *OrderStorage) *OrderHandler {
//	return &OrderHandler{
//		storage: storage,
//	}
//}
//
//func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
//	return h.storage.CreateOrder(req), nil
//}
//
//func (h *OrderHandler) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
//	order := h.storage.GetOrderById(params.OrderUUID)
//	if order == nil {
//		return &orderV1.NotFoundError{
//			Code:    http.StatusNotFound,
//			Message: "Order with UUID '" + params.OrderUUID.String() + "' not found",
//		}, nil
//	}
//
//	return order, nil
//}
//
//func (h *OrderHandler) PayOrderByUUID(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderByUUIDParams) (orderV1.PayOrderByUUIDRes, error) {
//	order := h.storage.PayOrderByUUID(ctx, params.OrderUUID, req.PaymentMethod)
//	if order == nil {
//		return &orderV1.NotFoundError{
//			Code:    http.StatusNotFound,
//			Message: "Order with UUID '" + params.OrderUUID.String() + "' not found",
//		}, nil
//	}
//
//	return order, nil
//}
//
//func (h *OrderHandler) CancelOrderByUUID(ctx context.Context, params orderV1.CancelOrderByUUIDParams) (orderV1.CancelOrderByUUIDRes, error) {
//	ok := h.storage.CancelOrderByUUID(params.OrderUUID)
//	if !ok {
//		return &orderV1.NotFoundError{
//			Code:    http.StatusNotFound,
//			Message: "Order with UUID '" + params.OrderUUID.String() + "' not found",
//		}, nil
//	}
//
//	return &orderV1.CancelOrderByUUIDNoContent{}, nil
//}

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to payment service: %v", err)
	}
	defer conn.Close()

	paymentClient := payment.NewClient(paymentV1.NewPaymentServiceClient(conn))
	if err != nil {
		log.Fatalf("Failed to create payment client: %v", err)
	}

	inventoryClient := inventory.NewClient(inventoryV1.NewInventoryServiceClient(conn))
	if err != nil {
		log.Fatalf("Failed to create inventory client: %v", err)
	}

	orderRepository := ordersRepo.NewRepository()
	orderService := ordersService.NewOrderService(orderRepository, paymentClient, inventoryClient)
	orderApi := ordersApi.NewAPI(orderService)

	orderHandler := 

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("Failed to create order server: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/api/v1", orderServer)

	server := &http.Server{
		Addr:    net.JoinHostPort("localhost", httpPort),
		Handler: r,
		ReadTimeout:  readHeaderTimeout,
	}

	go func() {
		log.Printf("Listening on port %s", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to listen on port %s: %v", httpPort, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("Failed to shutdown server gracefully: %v", err)
	}
	log.Println("Server gracefully stopped")
}
