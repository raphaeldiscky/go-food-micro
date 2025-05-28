package v1

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
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

// NewUpdateProduct is a constructor for the UpdateProduct.
func NewUpdateProduct(
	productID uuid.UUID,
	name string,
	description string,
	price float64,
) *UpdateProduct {
	command := &UpdateProduct{
		ProductID:   productID,
		Name:        name,
		Description: description,
		Price:       price,
		UpdatedAt:   time.Now(),
	}

	return command
}

// NewUpdateProductWithValidation is a constructor for the UpdateProduct with validation.
func NewUpdateProductWithValidation(
	productID uuid.UUID,
	name string,
	description string,
	price float64,
) (*UpdateProduct, error) {
	command := NewUpdateProduct(productID, name, description, price)
	err := command.Validate()

	return command, err
}

// Validate is a method that validates the update product command.
func (c *UpdateProduct) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ProductID, validation.Required),
		validation.Field(
			&c.Name,
			validation.Required,
			validation.Length(0, 255),
		),
		validation.Field(
			&c.Description,
			validation.Required,
			validation.Length(0, 5000),
		),
		validation.Field(&c.Price, validation.Required, validation.Min(0.0)),
		validation.Field(&c.UpdatedAt, validation.Required),
	)
	if err != nil {
		return customErrors.NewValidationErrorWrap(err, "validation error")
	}

	return nil
}
