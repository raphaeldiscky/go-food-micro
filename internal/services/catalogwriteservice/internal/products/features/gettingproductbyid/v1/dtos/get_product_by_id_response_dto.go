package dtos

import dtoV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1"

// https://echo.labstack.com/guide/response/

// GetProductByIDResponseDto is a struct that contains the get product by id response dto.
type GetProductByIDResponseDto struct {
	Product *dtoV1.ProductDto `json:"product"`
}
