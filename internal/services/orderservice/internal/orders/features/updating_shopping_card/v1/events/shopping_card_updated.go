package domainEvent

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/domain"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/value_objects"
)

type ShoppingCartUpdatedV1 struct {
	*domain.DomainEvent
	ShopItems []*value_objects.ShopItem `json:"shopItems" bson:"shopItems,omitempty"`
}

func NewShoppingCartUpdatedV1(shopItems []*value_objects.ShopItem) (*ShoppingCartUpdatedV1, error) {
	//if shopItems == nil {
	//	return nil, domainExceptions.ErrOrderShopItemsIsRequired
	//}

	eventData := ShoppingCartUpdatedV1{ShopItems: shopItems}

	return &eventData, nil
}
