package helpers

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/inventory/internal/model"
)

func CreatePart() model.Part {
	part := model.Part{
		PartId:        uuid.New(),
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.HipsterParagraph(),
		Price:         float32(gofakeit.Price(10, 1000)),
		StockQuantity: gofakeit.Int64(),
		Category:      model.CategoryEngine,
		Dimensions: model.Dimensions{
			Weight: gofakeit.Float64(),
			Length: gofakeit.Float64(),
			Width:  gofakeit.Float64()},
		Manufacturer: model.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL()},
		Tags:      []string{gofakeit.Word(), gofakeit.Word()},
		Metadata:  map[string]model.Value{"key": gofakeit.Word()},
		CreatedAt: gofakeit.Date(),
		UpdatedAt: gofakeit.Date(),
	}

	return part
}

func CreateParts(count int) []model.Part {
	parts := make([]model.Part, count)
	for i := 0; i < count; i++ {
		parts[i] = CreatePart()
	}

	return parts
}
