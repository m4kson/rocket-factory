package orders

import (
	"context"
	"testing"

	clientsMocks "github.com/m4kson/rocket-factory/order/internal/client/grpc/mocks"
	"github.com/m4kson/rocket-factory/order/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite

	orderRepository *mocks.OrderRepository

	paymentClient   *clientsMocks.PaymentClient
	inventoryClient *clientsMocks.InventoryClient

	ctx context.Context

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.paymentClient = clientsMocks.NewPaymentClient(s.T())
	s.inventoryClient = clientsMocks.NewInventoryClient(s.T())

	s.ctx = context.Background()

	s.service = NewOrderService(
		s.orderRepository,
		s.paymentClient,
		s.inventoryClient)
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
