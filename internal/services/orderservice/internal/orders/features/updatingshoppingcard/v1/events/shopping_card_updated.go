// Package domainevent contains the domain events for the shopping cart updated v1.
package domainevent

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/domain"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/valueobject"
)

// ShoppingCartUpdatedV1 is the event for the shopping cart updated v1.
type ShoppingCartUpdatedV1 struct {
	*domain.DomainEvent
	ShopItems []*valueobject.ShopItem `json:"shopItems" bson:"shopItems,omitempty"`
}

// NewShoppingCartUpdatedV1 creates a new shopping cart updated v1 event.
func NewShoppingCartUpdatedV1(shopItems []*valueobject.ShopItem) (*ShoppingCartUpdatedV1, error) {
	eventData := ShoppingCartUpdatedV1{ShopItems: shopItems}

	return &eventData, nil
}
