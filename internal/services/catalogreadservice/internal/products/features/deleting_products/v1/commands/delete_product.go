package commands

import (
	"github.com/go-ozzo/ozzo-validation/is"

	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"
)

// DeleteProduct is a struct that contains the delete product command.
type DeleteProduct struct {
	ProductID uuid.UUID
}

// NewDeleteProduct creates a new DeleteProduct.
func NewDeleteProduct(productID uuid.UUID) (*DeleteProduct, error) {
	delProduct := &DeleteProduct{ProductID: productID}
	if err := delProduct.Validate(); err != nil {
		return nil, err
	}

	return delProduct, nil
}

// Validate is a method that validates the delete product command.
func (p *DeleteProduct) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.ProductID, validation.Required, is.UUIDv4))
}
