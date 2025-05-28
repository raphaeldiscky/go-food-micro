package v1

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"
)

// CreateProduct is a struct that contains the create product command.
type CreateProduct struct {
	// we generate id ourselves because auto generate mongo string id column with type _id is not an uuid
	ID          string
	ProductID   string
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
}

// NewCreateProduct creates a new CreateProduct.
func NewCreateProduct(
	productId string,
	name string,
	description string,
	price float64,
	createdAt time.Time,
) (*CreateProduct, error) {
	command := &CreateProduct{
		ID:          uuid.NewV4().String(),
		ProductID:   productId,
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   createdAt,
	}
	if err := command.Validate(); err != nil {
		return nil, err
	}

	return command, nil
}

// Validate is a method that validates the create product command.
func (p *CreateProduct) Validate() error {
	return validation.ValidateStruct(p, validation.Field(&p.ID, validation.Required),
		validation.Field(&p.ProductID, validation.Required),
		validation.Field(&p.Name, validation.Required, validation.Length(3, 250)),
		validation.Field(&p.Description, validation.Required, validation.Length(3, 500)),
		validation.Field(&p.Price, validation.Required),
		validation.Field(&p.CreatedAt, validation.Required))
}
