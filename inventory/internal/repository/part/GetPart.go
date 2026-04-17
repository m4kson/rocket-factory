package part

import (
	"context"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/inventory/internal/model"
	"github.com/m4kson/rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) GetPart(ctx context.Context, partId uuid.UUID) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.parts[partId.String()]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	return converter.PartToModel(part), nil
}
