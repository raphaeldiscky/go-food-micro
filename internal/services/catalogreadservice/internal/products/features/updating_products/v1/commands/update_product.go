// Package commands contains the update product command.
package commands

import (
	"time"

	"github.com/go-ozzo/ozzo-validation/is"

	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"
)

// UpdateProduct is a struct that contains the update product command.
type UpdateProduct struct {
	ProductID   uuid.UUID
	Name        string
	Description string
	Price       float64
	UpdatedAt   time.Time
}

// NewUpdateProduct creates a new UpdateProduct.
func NewUpdateProduct(
	productId uuid.UUID,
	name string,
	description string,
	price float64,
) (*UpdateProduct, error) {
	product := &UpdateProduct{
		ProductID:   productId,
		Name:        name,
		Description: description,
		Price:       price,
		UpdatedAt:   time.Now(),
	}
	if err := product.Validate(); err != nil {
		return nil, err
	}

	return product, nil
}

// Validate is a method that validates the update product command.
func (p *UpdateProduct) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(&p.ProductID, validation.Required, is.UUIDv4),
		validation.Field(&p.Name, validation.Required, validation.Length(0, 255)),
		validation.Field(&p.Description, validation.Required, validation.Length(0, 5000)),
		validation.Field(&p.Price, validation.Required, validation.Min(0.0)),
		validation.Field(&p.UpdatedAt, validation.Required),
	)
}
