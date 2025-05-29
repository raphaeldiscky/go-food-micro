// Package domainexceptions contains the domain exceptions for the orderservice.
package domainexceptions

import (
	"emperror.dev/errors"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
)

// invalidDeliveryAddressError is the invalid delivery address error.
type invalidDeliveryAddressError struct {
	customErrors.BadRequestError
}

// NewInvalidDeliveryAddressError creates a new invalid delivery address error.
func NewInvalidDeliveryAddressError(message string) error {
	bad := customErrors.NewBadRequestError(message)
	customErr, ok := customErrors.GetCustomError(bad).(customErrors.BadRequestError)
	if !ok {
		return bad // Return original error if type assertion fails
	}

	br := &invalidDeliveryAddressError{
		BadRequestError: customErr,
	}

	return errors.WithStackIf(br)
}

// isInvalidAddress checks if the error is an invalid address error.
func (i *invalidDeliveryAddressError) isInvalidAddress() bool {
	return true
}

// IsInvalidDeliveryAddressError checks if the error is an invalid address error.
func IsInvalidDeliveryAddressError(err error) bool {
	var ia *invalidDeliveryAddressError
	if errors.As(err, &ia) {
		return ia.isInvalidAddress()
	}

	return false
}
