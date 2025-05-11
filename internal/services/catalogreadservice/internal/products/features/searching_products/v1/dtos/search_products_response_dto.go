package dtos

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/dto"
)

type SearchProductsResponseDto struct {
	Products *utils.ListResult[*dto.ProductDto]
}
