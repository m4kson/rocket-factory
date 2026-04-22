package converter

import (
	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/inventory/internal/model"
	repoModel "github.com/m4kson/rocket-factory/inventory/internal/repository/model"
)

func PartToRepoModel(part model.Part) repoModel.Part {
	return repoModel.Part{
		PartId:        part.PartId.String(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToRepoModel(part.Category),
		Dimensions:    DimensionsToRepoModel(part.Dimensions),
		Manufacturer:  ManufacturerToRepoModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToRepoModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func CategoryToRepoModel(category model.Category) repoModel.Category {
	switch category {
	case model.CategoryUnknown:
		return repoModel.CategoryUnknown
	case model.CategoryEngine:
		return repoModel.CategoryEngine
	case model.CategoryFuel:
		return repoModel.CategoryFuel
	case model.CategoryPorthole:
		return repoModel.CategoryPorthole
	case model.CategoryWing:
		return repoModel.CategoryWing
	default:
		return repoModel.CategoryUnknown
	}
}

func DimensionsToRepoModel(dimensions model.Dimensions) repoModel.Dimensions {
	return repoModel.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ManufacturerToRepoModel(manufacturer model.Manufacturer) repoModel.Manufacturer {
	return repoModel.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func MetadataToRepoModel(metadata map[string]model.Value) map[string]repoModel.Value {
	repoMetadata := make(map[string]repoModel.Value)
	for key, value := range metadata {
		repoMetadata[key] = ValueToRepoModel(value)
	}

	return repoMetadata
}

func ValueToRepoModel(value model.Value) repoModel.Value {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return v
	case float64:
		return v
	case bool:
		return v
	default:
		return nil
	}
}

func PartToModel(part repoModel.Part) model.Part {
	partId, err := uuid.Parse(part.PartId)
	if err != nil {
		partId = uuid.Nil
	}

	return model.Part{
		PartId:        partId,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToModel(part.Category),
		Dimensions:    DimensionsToModel(part.Dimensions),
		Manufacturer:  ManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func CategoryToModel(category repoModel.Category) model.Category {
	switch category {
	case repoModel.CategoryUnknown:
		return model.CategoryUnknown
	case repoModel.CategoryEngine:
		return model.CategoryEngine
	case repoModel.CategoryFuel:
		return model.CategoryFuel
	case repoModel.CategoryPorthole:
		return model.CategoryPorthole
	case repoModel.CategoryWing:
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}

func DimensionsToModel(dimensions repoModel.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ManufacturerToModel(manufacturer repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func MetadataToModel(metadata map[string]repoModel.Value) map[string]model.Value {
	modelMetadata := make(map[string]model.Value)
	for key, value := range metadata {
		modelMetadata[key] = ValueToModel(value)
	}

	return modelMetadata
}

func ValueToModel(value repoModel.Value) model.Value {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return v
	case float64:
		return v
	case bool:
		return v
	default:
		return nil
	}
}
