package dtos

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/serializer/json"

	uuid "github.com/satori/go.uuid"
)

// https://echo.labstack.com/guide/response/

// CreateProductResponseDto is a struct that contains the create product response dto.
type CreateProductResponseDto struct {
	ProductID uuid.UUID `json:"productID"`
}

// String is a method that returns the string representation of the create product response dto.
func (c *CreateProductResponseDto) String() string {
	return json.PrettyPrint(c)
}
