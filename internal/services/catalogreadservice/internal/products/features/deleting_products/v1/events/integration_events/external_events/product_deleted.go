package externalEvents

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
)

type ProductDeletedV1 struct {
	*types.Message
	ProductId string `json:"productId,omitempty"`
}
