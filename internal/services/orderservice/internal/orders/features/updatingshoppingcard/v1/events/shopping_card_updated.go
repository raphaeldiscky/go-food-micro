// Package domainevent contains the domain events for the shopping cart updated v1.
package domainevent

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/domain"

	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/valueobject"
)

// ShoppingCartUpdatedV1 is the event for the shopping cart updated v1.
type ShoppingCartUpdatedV1 struct {
	*domain.DomainEvent
	OrderId   uuid.UUID               `json:"orderId"   bson:"orderId,omitempty"`
	ShopItems []*valueobject.ShopItem `json:"shopItems" bson:"shopItems,omitempty"`
}

// NewShoppingCartUpdatedV1 creates a new shopping cart updated v1 event.
func NewShoppingCartUpdatedV1(
	orderId uuid.UUID,
	shopItems []*valueobject.ShopItem,
) (*ShoppingCartUpdatedV1, error) {
	eventData := &ShoppingCartUpdatedV1{
		OrderId:   orderId,
		ShopItems: shopItems,
	}
	eventData.DomainEvent = domain.NewDomainEvent(typeMapper.GetTypeName(eventData))

	return eventData, nil
}

// GetAggregateID returns the aggregate id.
func (e *ShoppingCartUpdatedV1) GetAggregateID() uuid.UUID {
	return e.OrderId
}
