// Package dtos contains the search products response dto.
package dtos

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/dto"
)

// SearchProductsResponseDto is a struct that contains the search products response dto.
type SearchProductsResponseDto struct {
	Products *utils.ListResult[*dto.ProductDto]
}
