// Package integrationevents contains the order created v1 event.
package integrationevents

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"

	uuid "github.com/satori/go.uuid"

	dtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/dtos/v1"
)

// OrderCreatedV1 is the order created v1 event.
type OrderCreatedV1 struct {
	*types.Message
	*dtosV1.OrderReadDto
}

// NewOrderCreatedV1 creates a new order created v1 event.
func NewOrderCreatedV1(orderReadDto *dtosV1.OrderReadDto) *OrderCreatedV1 {
	return &OrderCreatedV1{
		OrderReadDto: orderReadDto,
		Message:      types.NewMessage(uuid.NewV4().String()),
	}
}
