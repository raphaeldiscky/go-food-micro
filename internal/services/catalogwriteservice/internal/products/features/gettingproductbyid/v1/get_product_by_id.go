// Package v1 contains the get product by id query.
package v1

import (
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"

	validation "github.com/go-ozzo/ozzo-validation"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	uuid "github.com/satori/go.uuid"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

// GetProductByID is a struct that contains the get product by id query.
type GetProductByID struct {
	cqrs.Query
	ProductID uuid.UUID
}

// NewGetProductByID is a constructor for the GetProductByID.
func NewGetProductByID(productID uuid.UUID) *GetProductByID {
	query := &GetProductByID{
		Query:     cqrs.NewQueryByT[GetProductByID](),
		ProductID: productID,
	}

	return query
}

// NewGetProductByIDWithValidation is a constructor for the GetProductByID with validation.
func NewGetProductByIDWithValidation(productID uuid.UUID) (*GetProductByID, error) {
	query := NewGetProductByID(productID)
	err := query.Validate()

	return query, err
}

// Validate is a method that validates the get product by id query.
func (p *GetProductByID) Validate() error {
	err := validation.ValidateStruct(
		p,
		validation.Field(&p.ProductID, validation.Required, is.UUIDv4),
	)
	if err != nil {
		return customErrors.NewValidationErrorWrap(err, "validation error")
	}

	return nil
}
