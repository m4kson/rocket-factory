package orders

import (
	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestCancelOrderSuccess() {
	orderId := uuid.New()
	s.orderRepository.On("CancelOrderById", s.ctx, orderId).Return(nil)
	s.orderRepository.On("GetOrderById", s.ctx, orderId).Return(model.GetOrderResponse{}, nil)

	err := s.service.CancelOrderById(s.ctx, orderId)

	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderNotFound() {
	orderId := uuid.New()

	s.orderRepository.On("GetOrderById", s.ctx, orderId).Return(model.GetOrderResponse{}, model.ErrOrderNotFound)

	err := s.service.CancelOrderById(s.ctx, orderId)

	s.ErrorIs(err, model.ErrOrderNotFound)
}

func (s *ServiceSuite) TestCancelOrderAlreadyPaid() {
	orderId := uuid.New()

	s.orderRepository.On("CancelOrderById", s.ctx, orderId).Return(model.ErrOrderAlreadyPaid)
	s.orderRepository.On("GetOrderById", s.ctx, orderId).Return(model.GetOrderResponse{}, nil)

	err := s.service.CancelOrderById(s.ctx, orderId)

	s.ErrorIs(err, model.ErrOrderAlreadyPaid)
}
