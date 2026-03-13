package orders

import (
	"github.com/m4kson/rocket-factory/order/internal/repository"
	def "github.com/m4kson/rocket-factory/order/internal/service"

	grpc "github.com/m4kson/rocket-factory/order/internal/client/grpc"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository repository.OrderRepository

	paymentClient   grpc.PaymentClient
	inventoryClient grpc.InventoryClient
}

func NewOrderService(orderRepository repository.OrderRepository, paymentClient grpc.PaymentClient, inventoryClient grpc.InventoryClient) *service {
	return &service{
		orderRepository: orderRepository,
		paymentClient:   paymentClient,
		inventoryClient: inventoryClient,
	}
}
