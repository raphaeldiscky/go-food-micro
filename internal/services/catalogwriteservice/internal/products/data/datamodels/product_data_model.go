// Package datamodels contains the data models for the product.
package datamodels

import (
	"time"

	"gorm.io/gorm"

	json "github.com/goccy/go-json"
	uuid "github.com/satori/go.uuid"
)

// https://gorm.io/docs/conventions.html
// https://gorm.io/docs/models.html#gorm-Model

// ProductDataModel is a struct that contains the product data model.
type ProductDataModel struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time `gorm:"default:current_timestamp"`
	UpdatedAt   time.Time
	// for soft delete - https://gorm.io/docs/delete.html#Soft-Delete
	gorm.DeletedAt
}

// TableName overrides the table name used by ProductDataModel to `products` - https://gorm.io/docs/conventions.html#TableName
func (p *ProductDataModel) TableName() string {
	return "products"
}

// String is a method that returns the string representation of the product data model.
func (p *ProductDataModel) String() string {
	j, err := json.Marshal(p)
	if err != nil {
		return ""
	}

	return string(j)
}
