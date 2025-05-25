package externalEvents

import (
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
)

type ProductUpdatedV1 struct {
	*types.Message
	ProductId   string    `json:"productId,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (p *ProductUpdatedV1) GetMessageTypeName() string {
	return "ProductUpdatedV1"
}
