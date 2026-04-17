package v1

import "github.com/m4kson/rocket-factory/order/internal/service"

type api struct {
	orderService service.OrderService
}

func NewAPI(orderService service.OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
