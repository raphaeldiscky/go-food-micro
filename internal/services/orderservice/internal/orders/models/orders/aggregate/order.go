// Package aggregate contains the order aggregate.
package aggregate

// https://www.eventstore.com/blog/what-is-event-sourcing

import (
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/domain"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"

	json "github.com/goccy/go-json"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
	uuid "github.com/satori/go.uuid"

	dtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/dtos/v1"
	domainExceptions "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/exceptions/domain_exceptions"
	createOrderDomainEventsV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/creatingorder/v1/events/domain_events"
	updateOrderDomainEventsV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/updatingshoppingcard/v1/events"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/valueobject"
)

// Order is the order aggregate.
type Order struct {
	*models.EventSourcedAggregateRoot
	shopItems       []*valueobject.ShopItem
	accountEmail    string
	deliveryAddress string
	cancelReason    string
	deliveredTime   time.Time
	paid            bool
	submitted       bool
	completed       bool
	canceled        bool
	paymentId       uuid.UUID
	createdAt       time.Time
}

// NewEmptyAggregate creates a new empty aggregate.
func (o *Order) NewEmptyAggregate() {
	// http://arch-stable.blogspot.com/2012/05/golang-call-inherited-constructor.html
	base := models.NewEventSourcedAggregateRoot(typeMapper.GetFullTypeName(o), o.When)
	o.EventSourcedAggregateRoot = base
}

// NewOrder creates a new order aggregate.
func NewOrder(
	id uuid.UUID,
	shopItems []*valueobject.ShopItem,
	accountEmail, deliveryAddress string,
	deliveredTime time.Time,
	createdAt time.Time,
) (*Order, error) {
	order := &Order{}
	order.NewEmptyAggregate()
	order.SetId(id)

	if len(shopItems) == 0 {
		return nil, domainExceptions.NewOrderShopItemsRequiredError(
			"[Order_NewOrder] order items is required",
		)
	}

	itemsDto, err := mapper.Map[[]*dtosV1.ShopItemDto](shopItems)
	if err != nil {
		return nil, customErrors.NewDomainErrorWrap(
			err,
			"[Order_NewOrder.Map] error in the mapping []ShopItems to []ShopItemsDto",
		)
	}

	event, err := createOrderDomainEventsV1.NewOrderCreatedEventV1(
		id,
		itemsDto,
		accountEmail,
		deliveryAddress,
		deliveredTime,
		createdAt,
	)
	if err != nil {
		return nil, customErrors.NewDomainErrorWrap(
			err,
			"[Order_NewOrder.NewOrderCreatedEventV1] error in creating order created event",
		)
	}

	err = order.Apply(event, true)
	if err != nil {
		return nil, customErrors.NewDomainErrorWrap(
			err,
			"[Order_NewOrder.Apply] error in applying created event",
		)
	}

	return order, nil
}

// UpdateShoppingCard updates the shopping card.
func (o *Order) UpdateShoppingCard(shopItems []*valueobject.ShopItem) error {
	event, err := updateOrderDomainEventsV1.NewShoppingCartUpdatedV1(shopItems)
	if err != nil {
		return err
	}

	err = o.Apply(event, true)
	if err != nil {
		return err
	}

	return nil
}

// When handles the event.
func (o *Order) When(event domain.IDomainEvent) error {
	switch evt := event.(type) {
	case *createOrderDomainEventsV1.OrderCreatedV1:
		return o.onOrderCreated(evt)

	default:
		return errors.InvalidEventTypeError
	}
}

// onOrderCreated handles the order created event.
func (o *Order) onOrderCreated(evt *createOrderDomainEventsV1.OrderCreatedV1) error {
	items, err := mapper.Map[[]*valueobject.ShopItem](evt.ShopItems)
	if err != nil {
		return err
	}

	o.accountEmail = evt.AccountEmail
	o.shopItems = items
	o.deliveryAddress = evt.DeliveryAddress
	o.deliveredTime = evt.DeliveredTime
	o.createdAt = evt.CreatedAt
	o.SetId(evt.GetAggregateId()) // o.SetId(evt.ID)

	return nil
}

// ShopItems returns the shop items.
func (o *Order) ShopItems() []*valueobject.ShopItem {
	return o.shopItems
}

// PaymentId returns the payment id.
func (o *Order) PaymentId() uuid.UUID {
	return o.paymentId
}

// AccountEmail returns the account email.
func (o *Order) AccountEmail() string {
	return o.accountEmail
}

// DeliveryAddress returns the delivery address.
func (o *Order) DeliveryAddress() string {
	return o.deliveryAddress
}

// DeliveredTime returns the delivered time.
func (o *Order) DeliveredTime() time.Time {
	return o.deliveredTime
}

// CreatedAt returns the created at.
func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

// TotalPrice returns the total price.
func (o *Order) TotalPrice() float64 {
	return getShopItemsTotalPrice(o.shopItems)
}

// Paid returns the paid.
func (o *Order) Paid() bool {
	return o.paid
}

// Submitted returns the submitted.
func (o *Order) Submitted() bool {
	return o.submitted
}

// Completed returns the completed.
func (o *Order) Completed() bool {
	return o.completed
}

// Canceled returns the canceled.
func (o *Order) Canceled() bool {
	return o.canceled
}

// CancelReason returns the cancel reason.
func (o *Order) CancelReason() string {
	return o.cancelReason
}

// String returns the string representation of the order.
func (o *Order) String() string {
	j, _ := json.Marshal(o)

	return string(j)
}

// getShopItemsTotalPrice returns the total price of the shop items.
func getShopItemsTotalPrice(shopItems []*valueobject.ShopItem) float64 {
	var totalPrice float64 = 0
	for _, item := range shopItems {
		totalPrice += item.Price() * float64(item.Quantity())
	}

	return totalPrice
}
