// Package externalevents contains the product updated event.
package externalevents

import (
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
)

// ProductUpdatedV1 is a struct that contains the product updated event.
type ProductUpdatedV1 struct {
	*types.Message
	ProductID   string    `json:"productId,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// GetMessageTypeName is a method that returns the message type name.
func (p *ProductUpdatedV1) GetMessageTypeName() string {
	return "ProductUpdatedV1"
}
