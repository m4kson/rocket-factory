package model

import (
	"time"

	"github.com/google/uuid"
)

type Part struct {
	PartId        uuid.UUID
	Name          string
	Description   string
	Price         float32
	StockQuantity int64
	Category      Category
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
	Metadata      map[string]Value
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Category string

const (
	CategoryUnknown  Category = "unknown"
	CategoryEngine   Category = "engine"
	CategoryFuel     Category = "fuel"
	CategoryPorthole Category = "porthole"
	CategoryWing     Category = "wing"
)

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Value = any

type GetPartRequest struct {
	PartId uuid.UUID
}

type GetPartResponse struct {
	Part Part
}

type ListPartsRequest struct {
	Filter PartsFilter
}

type PartsFilter struct {
	Ids                   []uuid.UUID
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type ListPartsResponse struct {
	Parts []Part
}
