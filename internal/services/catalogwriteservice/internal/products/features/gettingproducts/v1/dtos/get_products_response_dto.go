package dtos

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
	dtoV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1"
)

// https://echo.labstack.com/guide/response/
type GetProductsResponseDto struct {
	Products *utils.ListResult[*dtoV1.ProductDto]
}
