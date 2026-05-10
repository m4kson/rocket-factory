package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Part struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	PartId        string           `bson:"part_id"`
	Name          string           `bson:"name"`
	Description   string           `bson:"description"`
	Price         float32          `bson:"price"`
	StockQuantity int64            `bson:"stock_quantity"`
	Category      Category         `bson:"category"`
	Dimensions    Dimensions       `bson:"dimensions"`
	Manufacturer  Manufacturer     `bson:"manufacturer"`
	Tags          []string         `bson:"tags"`
	Metadata      map[string]Value `bson:"metadata"`
	CreatedAt     time.Time        `bson:"created_at"`
	UpdatedAt     time.Time        `bson:"updated_at"`
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
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

type Value = any
