package dtos

import "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/dto"

type GetProductByIDResponseDto struct {
	Product *dto.ProductDto `json:"product"`
}
