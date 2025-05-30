// Package commands contains the commands for the submit order.
package commands

import (
	uuid "github.com/satori/go.uuid"
)

// SubmitOrder is the command for the submit order.
type SubmitOrder struct {
	OrderID uuid.UUID `validate:"required"`
}

// NewSubmitOrder creates a new submit order command.
func NewSubmitOrder(orderID uuid.UUID) *SubmitOrder {
	return &SubmitOrder{OrderID: orderID}
}
