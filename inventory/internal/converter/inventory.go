package converter

import (
	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
)

func GetPartRequestToModel(request *inventoryV1.GetPartRequest) uuid.UUID {
	var id uuid.UUID
	if request != nil {
		id = uuid.MustParse(request.GetPartUuid())
	}

	return id
}

func PartToProto(part model.Part) *inventoryV1.Part {
	return &inventoryV1.Part{
		Uuid:        part.PartId.String(),
		Name:        part.Name,
		Description: part.Description,
	}
}

func FilterToModel(filter *inventoryV1.PartsFilter) model.PartsFilter {
	if filter == nil {
		return model.PartsFilter{}
	}

	var ids []uuid.UUID
	if filter.Uuids != nil {
		for _, id := range filter.GetUuids() {
			ids = append(ids, uuid.MustParse(id))
		}
	}

	var names []string
	if filter.Names != nil {
		names = filter.GetNames()
	}

	var categories []model.Category
	if filter.Categories != nil {
		for _, category := range filter.GetCategories() {
			categories = append(categories, categoryToModel(category))
		}
	}

	var manufacturerCountries []string
	if filter.ManufacturerCountries != nil {
		manufacturerCountries = filter.GetManufacturerCountries()
	}

	var tags []string
	if filter.Tags != nil {
		tags = filter.GetTags()
	}

	return model.PartsFilter{
		Ids:                   ids,
		Names:                 names,
		Categories:            categories,
		ManufacturerCountries: manufacturerCountries,
		Tags:                  tags,
	}
}

func categoryToModel(c *inventoryV1.Category) model.Category {
	switch {
	case c.GetUnknown() != "":
		return model.CategoryUnknown
	case c.GetEngine() != "":
		return model.CategoryEngine
	case c.GetFuel() != "":
		return model.CategoryFuel
	case c.GetPorthole() != "":
		return model.CategoryPorthole
	case c.GetWing() != "":
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}
