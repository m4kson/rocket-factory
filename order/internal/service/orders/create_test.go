package orders

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
	"github.com/m4kson/rocket-factory/order/internal/service/orders/helpers"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestCreateOrderSuccess() {
	parts := helpers.CreateParts(3)

	serviceRequest := model.CreateOrderRequest{
		UserId:   uuid.New(),
		PartsIds: []uuid.UUID{parts[0].PartId, parts[1].PartId, parts[2].PartId},
	}

	response := model.CreateOrderRes{
		OrderId:    uuid.New(),
		TotalPrice: parts[0].Price + parts[1].Price + parts[2].Price,
	}

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Ids: serviceRequest.PartsIds}).Return(parts, nil)
	s.orderRepository.On("CreateOrder", s.ctx, mock.Anything).Return(response, nil)

	res, err := s.service.CreateOrder(s.ctx, serviceRequest)

	s.NoError(err)
	s.Equal(res, response)
}

func (s *ServiceSuite) TestCreateOrderInventoryError() {
	parts := helpers.CreateParts(3)

	serviceRequest := model.CreateOrderRequest{
		UserId:   uuid.New(),
		PartsIds: []uuid.UUID{parts[0].PartId, parts[1].PartId, parts[2].PartId},
	}

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Ids: serviceRequest.PartsIds}).Return(nil, gofakeit.Error())

	res, err := s.service.CreateOrder(s.ctx, serviceRequest)

	s.Error(err)
	s.Equal(res, model.CreateOrderRes{})
}

func (s *ServiceSuite) TestCreateOrderError() {
	parts := helpers.CreateParts(3)

	serviceRequest := model.CreateOrderRequest{
		UserId:   uuid.New(),
		PartsIds: []uuid.UUID{parts[0].PartId, parts[1].PartId, parts[2].PartId},
	}

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{Ids: serviceRequest.PartsIds}).Return(parts, nil)
	s.orderRepository.On("CreateOrder", s.ctx, mock.Anything).Return(model.CreateOrderRes{}, gofakeit.Error())

	res, err := s.service.CreateOrder(s.ctx, serviceRequest)

	s.Error(err)
	s.Equal(res, model.CreateOrderRes{})
}
