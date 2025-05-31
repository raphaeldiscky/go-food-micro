// Package externalevents contains the product created event.
package externalevents

import (
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
)

// ProductCreatedV1 is a struct that contains the product created event.
type ProductCreatedV1 struct {
	*types.Message
	ProductID   string    `json:"productID,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// GetMessageTypeName is a method that returns the message type name.
func (p *ProductCreatedV1) GetMessageTypeName() string {
	return "ProductCreatedV1"
}
