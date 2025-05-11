package integrationevents

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	dto "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1"

	uuid "github.com/satori/go.uuid"
)

type ProductUpdatedV1 struct {
	*types.Message
	*dto.ProductDto
}

func NewProductUpdatedV1(productDto *dto.ProductDto) *ProductUpdatedV1 {
	return &ProductUpdatedV1{
		Message:    types.NewMessage(uuid.NewV4().String()),
		ProductDto: productDto,
	}
}
