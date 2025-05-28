// Package v1 contains the delete product command.
package v1

import (
	"github.com/go-ozzo/ozzo-validation/is"

	validation "github.com/go-ozzo/ozzo-validation"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	uuid "github.com/satori/go.uuid"
)

// DeleteProduct is a struct that contains the delete product command.
type DeleteProduct struct {
	ProductID uuid.UUID
}

// NewDeleteProduct is a constructor for the DeleteProduct.
func NewDeleteProduct(productID uuid.UUID) *DeleteProduct {
	command := &DeleteProduct{ProductID: productID}

	return command
}

// NewDeleteProductWithValidation is a constructor for the DeleteProduct with validation.
func NewDeleteProductWithValidation(productID uuid.UUID) (*DeleteProduct, error) {
	command := NewDeleteProduct(productID)
	err := command.Validate()

	return command, err
}

// Validate is a method that validates the delete product command.
func (c *DeleteProduct) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ProductID, validation.Required),
		validation.Field(&c.ProductID, is.UUIDv4),
	)
	if err != nil {
		return customErrors.NewValidationErrorWrap(err, "validation error")
	}

	return nil
}
