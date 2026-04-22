package part

import (
	def "github.com/m4kson/rocket-factory/inventory/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	col *mongo.Collection
}

func NewPartRepository(col *mongo.Collection) *repository {
	return &repository{col: col}
}
