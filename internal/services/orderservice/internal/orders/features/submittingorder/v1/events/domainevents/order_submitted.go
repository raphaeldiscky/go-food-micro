// Package domainevents contains the domain events for the order submitted v1.
package domainevents

import (
	"fmt"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/domain"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
	uuid "github.com/satori/go.uuid"
)

// OrderSubmittedV1 is the event for the order submitted v1.
type OrderSubmittedV1 struct {
	*domain.DomainEvent
	OrderID uuid.UUID `json:"orderId" bson:"orderId,omitempty"`
}

// NewSubmitOrderV1 creates a new order submitted v1 event.
func NewSubmitOrderV1(orderID uuid.UUID) (*OrderSubmittedV1, error) {
	if orderID == uuid.Nil {
		return nil, customErrors.NewDomainError(fmt.Sprintf("orderId {%s} is invalid", orderID))
	}

	event := &OrderSubmittedV1{
		OrderID: orderID,
	}
	event.DomainEvent = domain.NewDomainEvent(typeMapper.GetTypeName(event))

	return event, nil
}

// GetAggregateID returns the aggregate id.
func (e *OrderSubmittedV1) GetAggregateID() uuid.UUID {
	return e.OrderID
}
