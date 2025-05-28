package integrationevents

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	dto "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1"

	uuid "github.com/satori/go.uuid"
)

// ProductUpdatedV1 is a struct that contains the product updated v1.
type ProductUpdatedV1 struct {
	*types.Message
	*dto.ProductDto
}

// NewProductUpdatedV1 is a constructor for the ProductUpdatedV1.
func NewProductUpdatedV1(productDto *dto.ProductDto) *ProductUpdatedV1 {
	return &ProductUpdatedV1{
		Message:    types.NewMessage(uuid.NewV4().String()),
		ProductDto: productDto,
	}
}
