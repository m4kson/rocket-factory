package part

import (
	"sync"

	def "github.com/m4kson/rocket-factory/inventory/internal/repository"
	repoModel "github.com/m4kson/rocket-factory/inventory/internal/repository/model"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	mu    sync.RWMutex
	parts map[string]repoModel.Part
}

func NewPartRepository() *repository {
	return &repository{
		parts: make(map[string]repoModel.Part),
	}
}
