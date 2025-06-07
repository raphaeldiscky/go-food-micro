// Package integrationevents contains the integration events for the product deleted v1.
package integrationevents

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"

	uuid "github.com/satori/go.uuid"
)

// ProductDeletedV1 is a struct that contains the product deleted v1.
type ProductDeletedV1 struct {
	*types.Message
	ProductID string `json:"productID,omitempty"`
}

// NewProductDeletedV1 is a constructor for the ProductDeletedV1.
func NewProductDeletedV1(productID string) *ProductDeletedV1 {
	return &ProductDeletedV1{ProductID: productID, Message: types.NewMessage(uuid.NewV4().String())}
}
