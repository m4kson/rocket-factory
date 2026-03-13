package orders

import (
	"sync"

	def "github.com/m4kson/rocket-factory/order/internal/repository"
	repoModel "github.com/m4kson/rocket-factory/order/internal/repository/model"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu     sync.RWMutex
	orders map[string]repoModel.Order
}

func NewRepository() *repository {
	return &repository{
		orders: make(map[string]repoModel.Order),
	}
}
