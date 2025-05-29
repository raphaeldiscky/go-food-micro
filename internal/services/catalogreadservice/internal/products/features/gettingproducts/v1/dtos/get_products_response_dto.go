// Package dtos contains the get products response dto.
package dtos

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/dto"
)

// GetProductsResponseDto is a struct that contains the get products response dto.
type GetProductsResponseDto struct {
	Products *utils.ListResult[*dto.ProductDto]
}
