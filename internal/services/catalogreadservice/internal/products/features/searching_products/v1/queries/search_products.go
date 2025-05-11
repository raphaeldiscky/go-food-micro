package queries

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	validation "github.com/go-ozzo/ozzo-validation"
)

type SearchProducts struct {
	SearchText string
	*utils.ListQuery
}

func (s *SearchProducts) Validate() error {
	return validation.ValidateStruct(s, validation.Field(&s.SearchText, validation.Required))
}
