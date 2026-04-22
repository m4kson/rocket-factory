package part

import (
	"context"
	"fmt"

	"github.com/m4kson/rocket-factory/inventory/internal/model"
	"github.com/m4kson/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/m4kson/rocket-factory/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *repository) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	mongoFilter := bson.D{}

	if filter.Ids != nil && len(filter.Ids) > 0 {
		for i := range filter.Ids {
			mongoFilter = append(mongoFilter, bson.E{Key: "part_id", Value: filter.Ids[i].String()})
		}
	}

	if filter.Names != nil && len(filter.Names) > 0 {
		for i := range filter.Names {
			mongoFilter = append(mongoFilter, bson.E{Key: "name", Value: filter.Names[i]})
		}
	}

	if filter.Categories != nil && len(filter.Categories) > 0 {
		for i := range filter.Categories {
			mongoFilter = append(mongoFilter, bson.E{Key: "category", Value: converter.CategoryToRepoModel(filter.Categories[i])})
		}
	}

	if filter.ManufacturerCountries != nil && len(filter.ManufacturerCountries) > 0 {
		for i := range filter.ManufacturerCountries {
			mongoFilter = append(mongoFilter, bson.E{Key: "manufacturer.country", Value: filter.ManufacturerCountries[i]})
		}
	}

	if filter.Tags != nil && len(filter.Tags) > 0 {
		for i := range filter.Tags {
			mongoFilter = append(mongoFilter, bson.E{Key: "tags", Value: filter.Tags[i]})
		}
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.col.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("repository.ListParts: %w", err)
	}
	defer cursor.Close(ctx)

	var rows []repoModel.Part
	if err = cursor.All(ctx, &rows); err != nil {
		return nil, fmt.Errorf("repository.ListParts: %w", err)
	}

	orders := make([]model.Part, 0, len(rows))
	for _, row := range rows {
		orders = append(orders, converter.PartToModel(row))
	}

	return orders, nil
}
