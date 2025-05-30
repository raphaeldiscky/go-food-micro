// Package models contains the product model.
package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Product is a struct that contains the product.
type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
