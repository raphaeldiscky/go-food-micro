// Package domainevents contains the domain events for the order submitted v1.
package domainevents

import (
	"fmt"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	uuid "github.com/satori/go.uuid"
)

// OrderSubmittedV1 is the event for the order submitted v1.
type OrderSubmittedV1 struct {
	OrderID uuid.UUID `json:"orderId" bson:"orderId,omitempty"`
}

// NewSubmitOrderV1 creates a new order submitted v1 event.
func NewSubmitOrderV1(orderID uuid.UUID) (*OrderSubmittedV1, error) {
	if orderID == uuid.Nil {
		return nil, customErrors.NewDomainError(fmt.Sprintf("orderId {%s} is invalid", orderID))
	}

	event := OrderSubmittedV1{OrderID: orderID}

	return &event, nil
}
