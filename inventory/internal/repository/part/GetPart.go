package part

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/inventory/internal/model"
	"github.com/m4kson/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/m4kson/rocket-factory/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) GetPart(ctx context.Context, partId uuid.UUID) (model.Part, error) {
	filter := bson.D{{Key: "part_id", Value: partId.String()}}

	var row repoModel.Part
	err := r.col.FindOne(ctx, filter).Decode(&row)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, model.ErrPartNotFound
		}

		return model.Part{}, fmt.Errorf("repository.GetPartById partid = %s : %w", partId.String(), err)
	}

	return converter.PartToModel(row), nil
}
