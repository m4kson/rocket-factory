package converter

import (
	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/order/internal/model"
	inventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartsFilterToProto(filter model.PartsFilter) *inventoryV1.PartsFilter {
	uuids := make([]string, 0, len(filter.Ids))
	for _, id := range filter.Ids {
		uuids = append(uuids, id.String())
	}

	categories := make([]*inventoryV1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, CategoryToProto(category))
	}

	return &inventoryV1.PartsFilter{
		Uuids:                 uuids,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func CategoryToProto(category model.Category) *inventoryV1.Category {
	switch category {
	case model.CategoryEngine:
		return &inventoryV1.Category{
			Category: &inventoryV1.Category_Engine{Engine: string(category)},
		}
	case model.CategoryFuel:
		return &inventoryV1.Category{
			Category: &inventoryV1.Category_Fuel{Fuel: string(category)},
		}
	case model.CategoryPorthole:
		return &inventoryV1.Category{
			Category: &inventoryV1.Category_Porthole{Porthole: string(category)},
		}
	case model.CategoryWing:
		return &inventoryV1.Category{
			Category: &inventoryV1.Category_Wing{Wing: string(category)},
		}
	case model.CategoryUnknown:
		fallthrough
	default:
		return &inventoryV1.Category{
			Category: &inventoryV1.Category_Unknown{Unknown: string(category)},
		}
	}
}

func PartToModel(part *inventoryV1.Part) model.Part {
	return model.Part{
		PartId:        uuid.MustParse(part.Uuid),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      ProtoCategoryToModel(part.Category),
		Dimensions: model.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Manufacturer: model.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags:      part.Tags,
		Metadata:  PartMetadataToModel(part.Metadata),
		CreatedAt: part.CreatedAt.AsTime(),
		UpdatedAt: part.UpdatedAt.AsTime(),
	}
}

func PartMetadataToModel(metadata map[string]*inventoryV1.Value) map[string]model.Value {
	result := make(map[string]model.Value, len(metadata))
	for key, value := range metadata {
		result[key] = value.Value
	}
	return result
}

func ProtoCategoryToModel(category *inventoryV1.Category) model.Category {
	switch category.Category.(type) {
	case *inventoryV1.Category_Engine:
		return model.CategoryEngine
	case *inventoryV1.Category_Fuel:
		return model.CategoryFuel
	case *inventoryV1.Category_Porthole:
		return model.CategoryPorthole
	case *inventoryV1.Category_Wing:
		return model.CategoryWing
	case *inventoryV1.Category_Unknown:
		return model.CategoryUnknown
	default:
		return model.CategoryUnknown
	}
}

func PartsListToModel(parts []*inventoryV1.Part) []model.Part {
	result := make([]model.Part, 0, len(parts))
	for _, part := range parts {
		result = append(result, PartToModel(part))
	}
	return result
}
