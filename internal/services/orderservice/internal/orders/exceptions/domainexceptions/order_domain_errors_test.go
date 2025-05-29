// Package domainexceptions contains the domain exceptions for the orderservice.
package domainexceptions

import (
	"testing"

	"github.com/stretchr/testify/assert"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
)

// TestOrderShopItemsRequiredError tests the order shop items required error.
func TestOrderShopItemsRequiredError(t *testing.T) {
	t.Parallel()
	err := NewOrderShopItemsRequiredError("order items required")
	assert.True(t, IsOrderShopItemsRequiredError(err))
}

// TestOrderNotFoundError tests the order not found error.
func TestOrderNotFoundError(t *testing.T) {
	t.Parallel()

	err := NewOrderNotFoundError(1)
	assert.True(t, IsOrderNotFoundError(err))
}

// TestInvalidDeliveryAddressError tests the invalid delivery address error.
func TestInvalidDeliveryAddressError(t *testing.T) {
	t.Parallel()

	err := NewInvalidDeliveryAddressError("address is not valid")
	assert.True(t, IsInvalidDeliveryAddressError(err))
}

// TestIsNotInvalidDeliveryAddressError tests the is not invalid delivery address error.
func TestIsNotInvalidDeliveryAddressError(
	t *testing.T,
) {
	t.Parallel()

	err := customErrors.NewBadRequestError("address is not valid")
	assert.False(t, IsInvalidDeliveryAddressError(err))
}

// TestInvalidEmailAddressError tests the invalid email address error.
func TestInvalidEmailAddressError(t *testing.T) {
	t.Parallel()

	err := NewInvalidEmailAddressError("email address is not valid")
	assert.True(t, IsInvalidEmailAddressError(err))
}

// TestIsNotInvalidEmailAddressError tests the is not invalid email address error.
func TestIsNotInvalidEmailAddressError(t *testing.T) {
	t.Parallel()

	err := customErrors.NewBadRequestError("email address is not valid")
	assert.False(t, IsInvalidEmailAddressError(err))
}
