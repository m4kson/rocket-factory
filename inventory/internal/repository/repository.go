package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/inventory/internal/model"
)

type PartRepository interface {
	GetPart(ctx context.Context, partId uuid.UUID) (model.Part, error)
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
