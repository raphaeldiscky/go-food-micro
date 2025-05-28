package queries

import (
	"github.com/go-ozzo/ozzo-validation/is"

	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"
)

// GetProductByID is a struct that contains the get product by id query.
type GetProductByID struct {
	ID uuid.UUID
}

// NewGetProductByID creates a new GetProductByID.
func NewGetProductByID(id uuid.UUID) (*GetProductByID, error) {
	product := &GetProductByID{ID: id}
	if err := product.Validate(); err != nil {
		return nil, err
	}

	return product, nil
}

// Validate is a method that validates the get product by id query.
func (p *GetProductByID) Validate() error {
	return validation.ValidateStruct(p, validation.Field(&p.ID, validation.Required, is.UUIDv4))
}
