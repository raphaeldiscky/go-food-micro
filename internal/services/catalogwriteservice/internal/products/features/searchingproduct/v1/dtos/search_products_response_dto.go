package dtos

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	dtoV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1"
)

// SearchProductsResponseDto is a struct that contains the search products response dto.
type SearchProductsResponseDto struct {
	Products *utils.ListResult[*dtoV1.ProductDto]
}
