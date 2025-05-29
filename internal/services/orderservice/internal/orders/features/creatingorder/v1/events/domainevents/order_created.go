// Package domainevents contains the order created v1 event.
package domainevents

import (
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/domain"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
	uuid "github.com/satori/go.uuid"

	dtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/dtos/v1"
	domainExceptions "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/exceptions/domainexceptions"
)

// OrderCreatedV1 is the order created v1 event.
type OrderCreatedV1 struct {
	*domain.DomainEvent
	OrderID         uuid.UUID             `json:"order_id"`
	ShopItems       []*dtosV1.ShopItemDto `json:"shopItems"       bson:"shopItems,omitempty"`
	AccountEmail    string                `json:"accountEmail"    bson:"accountEmail,omitempty"`
	DeliveryAddress string                `json:"deliveryAddress" bson:"deliveryAddress,omitempty"`
	CreatedAt       time.Time             `json:"createdAt"       bson:"createdAt,omitempty"`
	DeliveredTime   time.Time             `json:"deliveredTime"   bson:"deliveredTime,omitempty"`
}

// NewOrderCreatedEventV1 creates a new order created v1 event.
func NewOrderCreatedEventV1(
	orderID uuid.UUID,
	shopItems []*dtosV1.ShopItemDto,
	accountEmail, deliveryAddress string,
	deliveredTime time.Time,
	createdAt time.Time,
) (*OrderCreatedV1, error) {
	if len(shopItems) == 0 {
		return nil, domainExceptions.NewOrderShopItemsRequiredError("shopItems is required")
	}

	if deliveryAddress == "" {
		return nil, domainExceptions.NewInvalidDeliveryAddressError("deliveryAddress is invalid")
	}

	if accountEmail == "" {
		return nil, domainExceptions.NewInvalidEmailAddressError("accountEmail is invalid")
	}

	if createdAt.IsZero() {
		return nil, customErrors.NewDomainError("createdAt can't be zero")
	}

	if deliveredTime.IsZero() {
		return nil, customErrors.NewDomainError("deliveredTime can't be zero")
	}

	eventData := &OrderCreatedV1{
		ShopItems:       shopItems,
		OrderID:         orderID,
		AccountEmail:    accountEmail,
		DeliveryAddress: deliveryAddress,
		CreatedAt:       createdAt,
		DeliveredTime:   deliveredTime,
	}

	eventData.DomainEvent = domain.NewDomainEvent(typeMapper.GetTypeName(eventData))

	return eventData, nil
}

// GetAggregateID returns the aggregate id.
func (e *OrderCreatedV1) GetAggregateID() uuid.UUID {
	return e.OrderID
}
