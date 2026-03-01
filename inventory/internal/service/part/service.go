package part

import (
	"github.com/m4kson/rocket-factory/inventory/internal/repository"
	def "github.com/m4kson/rocket-factory/inventory/internal/service"
)

var _ def.PartService = (*service)(nil)

type service struct {
	partRepository repository.PartRepository
}

func NewPartService(partRepository repository.PartRepository) *service {
	return &service{
		partRepository: partRepository,
	}
}
