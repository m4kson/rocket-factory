package part

import (
	"context"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, partId uuid.UUID) (model.Part, error) {
	part, err := s.partRepository.GetPart(ctx, partId)
	if err != nil {
		return model.Part{}, err
	}

	return part, nil
}
