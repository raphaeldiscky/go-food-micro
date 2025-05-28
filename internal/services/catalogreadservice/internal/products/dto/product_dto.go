package dto

import (
	"time"
)

// ProductDto is a struct that contains the product dto.
type ProductDto struct {
	ID          string    `json:"id"`
	ProductID   string    `json:"productId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
