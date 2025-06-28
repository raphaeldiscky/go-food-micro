// Package externalevents contains the product deleted event.
package externalevents

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
)

// ProductDeletedV1 is a struct that contains the product deleted event.
type ProductDeletedV1 struct {
	*types.Message
	ProductID string `json:"productID,omitempty"`
}

// GetMessageTypeName returns the message type name.
func (p *ProductDeletedV1) GetMessageTypeName() string {
	return "ProductDeletedV1"
}
