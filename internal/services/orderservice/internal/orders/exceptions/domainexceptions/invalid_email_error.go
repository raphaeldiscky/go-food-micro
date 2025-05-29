// Package domainexceptions contains the domain exceptions for the orderservice.
package domainexceptions

import (
	"emperror.dev/errors"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
)

// invalidEmailAddressError is the invalid email address error.
type invalidEmailAddressError struct {
	customErrors.BadRequestError
}

// NewInvalidEmailAddressError creates a new invalid email address error.
func NewInvalidEmailAddressError(message string) error {
	bad := customErrors.NewBadRequestError(message)
	customErr, ok := customErrors.GetCustomError(bad).(customErrors.BadRequestError)
	if !ok {
		return bad // Return original error if type assertion fails
	}

	br := &invalidEmailAddressError{
		BadRequestError: customErr,
	}

	return errors.WithStackIf(br)
}

// isInvalidEmailAddressError checks if the error is an invalid email address error.
func (i *invalidEmailAddressError) isInvalidEmailAddressError() bool {
	return true
}

// IsInvalidEmailAddressError checks if the error is an invalid email address error.
func IsInvalidEmailAddressError(err error) bool {
	var ie *invalidEmailAddressError

	if errors.As(err, &ie) {
		return ie.isInvalidEmailAddressError()
	}

	return false
}
