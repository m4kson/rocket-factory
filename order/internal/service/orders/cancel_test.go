package orders

import (
	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestCancelOrderSuccess() {
	orderId := uuid.New()
	s.orderRepository.On("CancelOrderById", s.ctx, orderId).Return(nil)

	err := s.service.CancelOrderById(s.ctx, orderId)

	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderNotFound() {
	orderId := uuid.New()
	s.orderRepository.On("CancelOrderById", s.ctx, orderId).Return(model.ErrOrderNotFound)

	err := s.service.CancelOrderById(s.ctx, orderId)

	s.ErrorIs(err, model.ErrOrderNotFound)
}

func (s *ServiceSuite) TestCancelOrderAlreadyPaid() {
	orderId := uuid.New()

	s.orderRepository.On("CancelOrderById", s.ctx, orderId).Return(model.ErrOrderAlreadyPaid)

	err := s.service.CancelOrderById(s.ctx, orderId)

	s.ErrorIs(err, model.ErrOrderAlreadyPaid)
}
