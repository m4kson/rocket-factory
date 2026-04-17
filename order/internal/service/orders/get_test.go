package orders

import (
	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestGetOrderByIdSuccess() {
	order := model.GetOrderResponse{
		OrderId:       uuid.New(),
		UserId:        uuid.New(),
		PartsIds:      []uuid.UUID{uuid.New(), uuid.New()},
		TotalPrice:    100.0,
		TransactionId: nil,
		PaymentMethod: model.PaymentMethodCARD,
		Status:        model.OrderStatusPENDINGPAYMENT,
	}

	s.orderRepository.On("GetOrderById", s.ctx, order.OrderId).Return(order, nil)

	res, err := s.service.GetOrderById(s.ctx, order.OrderId)

	s.NoError(err)
	s.Equal(order, res)
}
