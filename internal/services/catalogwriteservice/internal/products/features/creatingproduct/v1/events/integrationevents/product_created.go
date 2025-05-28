package integrationevents

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	dtoV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1"

	uuid "github.com/satori/go.uuid"
)

// ProductCreatedV1 is a struct that contains the product created v1.
type ProductCreatedV1 struct {
	*types.Message
	*dtoV1.ProductDto
}

// NewProductCreatedV1 is a constructor for the ProductCreatedV1.
func NewProductCreatedV1(productDto *dtoV1.ProductDto) *ProductCreatedV1 {
	return &ProductCreatedV1{
		ProductDto: productDto,
		Message:    types.NewMessage(uuid.NewV4().String()),
	}
}
