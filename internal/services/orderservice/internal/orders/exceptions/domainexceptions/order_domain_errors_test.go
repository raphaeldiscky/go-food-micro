// Package domainexceptions contains the domain exceptions for the orderservice.
package domainexceptions

import (
	"testing"

	"github.com/stretchr/testify/assert"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
)

func Test_Order_Shop_Items_Required_Error(t *testing.T) {
	t.Parallel()
	err := NewOrderShopItemsRequiredError("order items required")
	assert.True(t, IsOrderShopItemsRequiredError(err))
}

func Test_Order_Not_Found_Error(t *testing.T) {
	t.Parallel()

	err := NewOrderNotFoundError(1)
	assert.True(t, IsOrderNotFoundError(err))
}

func Test_Invalid_Delivery_Address_Error(t *testing.T) {
	t.Parallel()

	err := NewInvalidDeliveryAddressError("address is not valid")
	assert.True(t, IsInvalidDeliveryAddressError(err))
}

func Test_Is_Not_Invalid_Delivery_Address_Error(
	t *testing.T,
) {
	t.Parallel()

	err := customErrors.NewBadRequestError("address is not valid")
	assert.False(t, IsInvalidDeliveryAddressError(err))
}

func Test_InvalidEmail_Address_Error(t *testing.T) {
	t.Parallel()

	err := NewInvalidEmailAddressError("email address is not valid")
	assert.True(t, IsInvalidEmailAddressError(err))
}

func Test_Is_Not_InvalidEmail_Address_Error(t *testing.T) {
	t.Parallel()

	err := customErrors.NewBadRequestError("email address is not valid")
	assert.False(t, IsInvalidEmailAddressError(err))
}
