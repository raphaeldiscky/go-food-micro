// Package dtosv1 contains the order read dto.
package dtosv1

import "time"

// OrderReadDto is the read dto for the order.
type OrderReadDto struct {
	ID              string             `json:"id"`
	OrderID         string             `json:"orderId"`
	ShopItems       []*ShopItemReadDto `json:"shopItems"`
	AccountEmail    string             `json:"accountEmail"`
	DeliveryAddress string             `json:"deliveryAddress"`
	CancelReason    string             `json:"cancelReason"`
	TotalPrice      float64            `json:"totalPrice"`
	DeliveredTime   time.Time          `json:"deliveredTime"`
	Paid            bool               `json:"paid"`
	Submitted       bool               `json:"submitted"`
	Completed       bool               `json:"completed"`
	Canceled        bool               `json:"canceled"`
	PaymentID       string             `json:"paymentID"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
}
