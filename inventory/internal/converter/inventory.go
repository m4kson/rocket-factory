package converter

import (
	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		Uuid:          part.PartId.String(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      categoryToProto(part.Category),
		Dimensions:    dimensionsToProto(part.Dimensions),
		Manufacturer:  manufacturerToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      metadataToProto(part.Metadata),
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     timestamppb.New(part.UpdatedAt),
	}
}

func metadataToProto(metadata map[string]model.Value) map[string]*inventoryV1.Value {
	protoMetadata := make(map[string]*inventoryV1.Value)
	for key, value := range metadata {
		protoMetadata[key] = valueToProto(value)
	}

	return protoMetadata
}

func valueToProto(value model.Value) *inventoryV1.Value {
	switch v := value.(type) {
	case string:
		return &inventoryV1.Value{
			Value: &inventoryV1.Value_StringValue{StringValue: v},
		}
	case int64:
		return &inventoryV1.Value{
			Value: &inventoryV1.Value_Int64Value{Int64Value: v},
		}
	case int:
		return &inventoryV1.Value{
			Value: &inventoryV1.Value_Int64Value{Int64Value: int64(v)},
		}
	case float64:
		return &inventoryV1.Value{
			Value: &inventoryV1.Value_DoubleValue{DoubleValue: v},
		}
	case bool:
		return &inventoryV1.Value{
			Value: &inventoryV1.Value_BoolValue{BoolValue: v},
		}
	default:
		return nil
	}
}

func manufacturerToProto(manufacturer model.Manufacturer) *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func dimensionsToProto(dimensions model.Dimensions) *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
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
	if c == nil {
		return model.CategoryUnknown
	}

	switch c.Category.(type) {
	case *inventoryV1.Category_Engine:
		return model.CategoryEngine
	case *inventoryV1.Category_Fuel:
		return model.CategoryFuel
	case *inventoryV1.Category_Porthole:
		return model.CategoryPorthole
	case *inventoryV1.Category_Wing:
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}

func categoryToProto(category model.Category) *inventoryV1.Category {
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
	default:
		return &inventoryV1.Category{
			Category: &inventoryV1.Category_Unknown{Unknown: string(category)},
		}
	}
}
