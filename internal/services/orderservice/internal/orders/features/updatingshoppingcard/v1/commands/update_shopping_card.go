// Package commands contains the commands for the update shopping cart.
package commands

import (
	uuid "github.com/satori/go.uuid"

	dtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/dtos/v1"
)

// UpdateShoppingCart is the command for the update shopping cart.
type UpdateShoppingCart struct {
	OrderID   uuid.UUID             `validate:"required"`
	ShopItems []*dtosV1.ShopItemDto `validate:"required"`
}

// NewUpdateShoppingCart creates a new update shopping cart command.
func NewUpdateShoppingCart(orderID uuid.UUID, shopItems []*dtosV1.ShopItemDto) *UpdateShoppingCart {
	return &UpdateShoppingCart{OrderID: orderID, ShopItems: shopItems}
}
