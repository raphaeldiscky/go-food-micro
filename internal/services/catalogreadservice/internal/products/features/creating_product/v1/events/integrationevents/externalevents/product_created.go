package externalEvents

import (
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
)

type ProductCreatedV1 struct {
	*types.Message
	ProductId   string    `json:"productId,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (p *ProductCreatedV1) GetMessageTypeName() string {
	return "ProductCreatedV1"
}
