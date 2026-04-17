package orders

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestPaySuccess() {
	orderId := uuid.New()
	paymentMethod := model.PaymentMethodSBP
	userId := uuid.New()
	transactionId := uuid.New()

	payOrderRequest := model.PayOrderRequest{
		UserId:        userId.String(),
		PaymentMethod: paymentMethod,
		OrderId:       orderId.String(),
	}

	payOrderResponse := model.PayOrderResponse{
		TransactionId: transactionId.String(),
	}

	payOrderRes := model.PayOrderRes{
		TransactionId: transactionId,
	}

	s.paymentClient.On("PayOrder", s.ctx, payOrderRequest).Return(payOrderResponse, nil).Once()
	s.orderRepository.On("PayOrderById", s.ctx, orderId, paymentMethod, uuid.MustParse(payOrderResponse.TransactionId)).Return(payOrderRes, nil)

	response, err := s.service.PayOrderById(s.ctx, orderId, paymentMethod, userId)
	s.NoError(err)
	s.Equal(payOrderRes, response)
}

func (s *ServiceSuite) TestPayPaymentClientError() {
	orderId := uuid.New()
	paymentMethod := model.PaymentMethodSBP
	userId := uuid.New()

	payOrderRequest := model.PayOrderRequest{
		UserId:        userId.String(),
		PaymentMethod: paymentMethod,
		OrderId:       orderId.String(),
	}

	s.paymentClient.On("PayOrder", s.ctx, payOrderRequest).Return(model.PayOrderResponse{}, gofakeit.Error()).Once()

	response, err := s.service.PayOrderById(s.ctx, orderId, paymentMethod, userId)

	s.Error(err)
	s.Equal(response, model.PayOrderRes{})
}

func (s *ServiceSuite) TestPayRepositoryError() {
	orderId := uuid.New()
	paymentMethod := model.PaymentMethodSBP
	userId := uuid.New()
	transactionId := uuid.New()

	payOrderRequest := model.PayOrderRequest{
		UserId:        userId.String(),
		PaymentMethod: paymentMethod,
		OrderId:       orderId.String(),
	}

	payOrderResponse := model.PayOrderResponse{
		TransactionId: transactionId.String(),
	}

	s.paymentClient.On("PayOrder", s.ctx, payOrderRequest).Return(payOrderResponse, nil).Once()
	s.orderRepository.On("PayOrderById", s.ctx, orderId, paymentMethod, uuid.MustParse(payOrderResponse.TransactionId)).Return(model.PayOrderRes{}, gofakeit.Error())

	response, err := s.service.PayOrderById(s.ctx, orderId, paymentMethod, userId)
	s.Error(err)
	s.Equal(response, model.PayOrderRes{})
}
