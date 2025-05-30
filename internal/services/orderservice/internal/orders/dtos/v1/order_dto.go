// Package dtosv1 contains the order dto.
package dtosv1

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// OrderDto is the dto for the order.
type OrderDto struct {
	ID              uuid.UUID      `json:"id"`
	ShopItems       []*ShopItemDto `json:"shopItems"`
	AccountEmail    string         `json:"accountEmail"`
	DeliveryAddress string         `json:"deliveryAddress"`
	CancelReason    string         `json:"cancelReason"`
	TotalPrice      float64        `json:"totalPrice"`
	DeliveredTime   time.Time      `json:"deliveredTime"`
	Paid            bool           `json:"paid"`
	Submitted       bool           `json:"submitted"`
	Completed       bool           `json:"completed"`
	Canceled        bool           `json:"canceled"`
	PaymentID       uuid.UUID      `json:"paymentID"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	OriginalVersion int64          `json:"originalVersion"`
}
