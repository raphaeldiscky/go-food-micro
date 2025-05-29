// Package domainexceptions contains the domain exceptions for the orderservice.
package domainexceptions

import (
	"fmt"

	"emperror.dev/errors"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
)

// orderNotFoundError is the order not found error.
type orderNotFoundError struct {
	customErrors.NotFoundError
}

// NewOrderNotFoundError creates a new order not found error.
func NewOrderNotFoundError(id int) error {
	notFound := customErrors.NewNotFoundError(
		fmt.Sprintf("order with id %d not found", id),
	)
	customErr, ok := customErrors.GetCustomError(notFound).(customErrors.NotFoundError)
	if !ok {
		return notFound // Return original error if type assertion fails
	}

	br := &orderNotFoundError{
		NotFoundError: customErr,
	}

	return errors.WithStackIf(br)
}

// isorderNotFoundError checks if the error is an order not found error.
func (i *orderNotFoundError) isorderNotFoundError() bool {
	return true
}

// IsOrderNotFoundError checks if the error is an order not found error.
func IsOrderNotFoundError(err error) bool {
	var os *orderNotFoundError
	if errors.As(err, &os) {
		return os.isorderNotFoundError()
	}

	return false
}
