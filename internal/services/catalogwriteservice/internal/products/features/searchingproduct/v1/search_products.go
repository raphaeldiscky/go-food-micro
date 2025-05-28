package v1

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	validation "github.com/go-ozzo/ozzo-validation"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
)

// SearchProducts is a struct that contains the search products query.
type SearchProducts struct {
	SearchText string
	*utils.ListQuery
}

// NewSearchProducts is a constructor for the SearchProducts.
func NewSearchProducts(searchText string, query *utils.ListQuery) *SearchProducts {
	searchProductQuery := &SearchProducts{
		SearchText: searchText,
		ListQuery:  query,
	}

	return searchProductQuery
}

// NewSearchProductsWithValidation is a constructor for the SearchProducts with validation.
func NewSearchProductsWithValidation(
	searchText string,
	query *utils.ListQuery,
) (*SearchProducts, error) {
	searchProductQuery := NewSearchProducts(searchText, query)

	err := searchProductQuery.Validate()

	return searchProductQuery, err
}

// Validate is a method that validates the search products query.
func (p *SearchProducts) Validate() error {
	err := validation.ValidateStruct(p, validation.Field(&p.SearchText, validation.Required))
	if err != nil {
		return customErrors.NewValidationErrorWrap(err, "validation error")
	}

	return nil
}
