package queries

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	validation "github.com/go-ozzo/ozzo-validation"
)

// SearchProducts is a struct that contains the search products query.
type SearchProducts struct {
	SearchText string
	*utils.ListQuery
}

// Validate is a method that validates the search products query.
func (s *SearchProducts) Validate() error {
	return validation.ValidateStruct(s, validation.Field(&s.SearchText, validation.Required))
}
