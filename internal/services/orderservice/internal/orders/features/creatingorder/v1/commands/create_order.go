// Package createordercommandv1 contains the create order command.
package createordercommandv1

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"

	dtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/dtos/v1"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

// CreateOrder is the create order command.
type CreateOrder struct {
	OrderID         uuid.UUID
	ShopItems       []*dtosV1.ShopItemDto
	AccountEmail    string
	DeliveryAddress string
	DeliveryTime    time.Time
	CreatedAt       time.Time
}

// NewCreateOrder creates a new create order command.
func NewCreateOrder(
	shopItems []*dtosV1.ShopItemDto,
	accountEmail, deliveryAddress string,
	deliveryTime time.Time,
) (*CreateOrder, error) {
	command := &CreateOrder{
		OrderID:         uuid.NewV4(),
		ShopItems:       shopItems,
		AccountEmail:    accountEmail,
		DeliveryAddress: deliveryAddress,
		DeliveryTime:    deliveryTime,
		CreatedAt:       time.Now(),
	}

	err := command.Validate()
	if err != nil {
		return nil, err
	}

	return command, nil
}

// Validate validates the create order command.
func (c *CreateOrder) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.OrderID, validation.Required),
		validation.Field(&c.ShopItems, validation.Required),
		validation.Field(&c.AccountEmail, validation.Required),
		validation.Field(&c.DeliveryAddress, validation.Required),
		validation.Field(&c.DeliveryTime, validation.Required),
		validation.Field(&c.CreatedAt, validation.Required),
	)
}
