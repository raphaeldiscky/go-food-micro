// Package domainexceptions contains the domain exceptions for the orderservice.
package domainexceptions

import (
	"emperror.dev/errors"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
)

// orderShopItemsRequiredError is the order shop items required error.
type orderShopItemsRequiredError struct {
	customErrors.BadRequestError
}

// NewOrderShopItemsRequiredError creates a new order shop items required error.
func NewOrderShopItemsRequiredError(message string) error {
	bad := customErrors.NewBadRequestError(message)
	customErr := customErrors.GetCustomError(bad).(customErrors.BadRequestError)
	br := &orderShopItemsRequiredError{
		BadRequestError: customErr,
	}

	return errors.WithStackIf(br)
}

// isOrderShopItemsRequiredError checks if the error is an order shop items required error.
func (i *orderShopItemsRequiredError) isOrderShopItemsRequiredError() bool {
	return true
}

// IsOrderShopItemsRequiredError checks if the error is an order shop items required error.
func IsOrderShopItemsRequiredError(err error) bool {
	var os *orderShopItemsRequiredError
	if errors.As(err, &os) {
		return os.isOrderShopItemsRequiredError()
	}

	return false
}
