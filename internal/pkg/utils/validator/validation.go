// Package validator provides a validator.
package validator

import (
	"errors"

	"github.com/go-playground/validator"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate attempts to validate the lead's values.
func Validate(input interface{}) error {
	if err := validate.Struct(input); err != nil {
		// this check ensures there wasn't an error
		// with the validation process itself
		invalidValidationError := &validator.InvalidValidationError{}
		if errors.As(err, &invalidValidationError) {
			return err
		}

		return err
	}

	return nil
}
