// Package dtos contains the get product by id response dto.
package dtos

import "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/dto"

// GetProductByIDResponseDto is a struct that contains the get product by id response dto.
type GetProductByIDResponseDto struct {
	Product *dto.ProductDto `json:"product"`
}
