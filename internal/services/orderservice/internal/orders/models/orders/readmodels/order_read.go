// Package readmodels contains the order read model.
package readmodels

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// OrderReadModel is the read model for the order.
type OrderReadModel struct {
	// we generate id ourself because auto generate mongo string id column with type _id is not an uuid
	ID              string               `json:"id"                        bson:"_id,omitempty"` // https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/write-operations/insert/#the-_id-field
	OrderID         string               `json:"orderId"                   bson:"orderId,omitempty"`
	ShopItems       []*ShopItemReadModel `json:"shopItems,omitempty"       bson:"shopItems,omitempty"`
	AccountEmail    string               `json:"accountEmail,omitempty"    bson:"accountEmail,omitempty"`
	DeliveryAddress string               `json:"deliveryAddress,omitempty" bson:"deliveryAddress,omitempty"`
	CancelReason    string               `json:"cancelReason,omitempty"    bson:"cancelReason,omitempty"`
	TotalPrice      float64              `json:"totalPrice,omitempty"      bson:"totalPrice,omitempty"`
	DeliveredTime   time.Time            `json:"deliveredTime,omitempty"   bson:"deliveredTime,omitempty"`
	Paid            bool                 `json:"paid,omitempty"            bson:"paid,omitempty"`
	Submitted       bool                 `json:"submitted,omitempty"       bson:"submitted,omitempty"`
	Completed       bool                 `json:"completed,omitempty"       bson:"completed,omitempty"`
	Canceled        bool                 `json:"canceled,omitempty"        bson:"canceled,omitempty"`
	PaymentID       string               `json:"paymentID"                 bson:"paymentID,omitempty"`
	CreatedAt       time.Time            `json:"createdAt,omitempty"       bson:"createdAt,omitempty"`
	UpdatedAt       time.Time            `json:"updatedAt,omitempty"       bson:"updatedAt,omitempty"`
}

// NewOrderReadModel creates a new order read model.
func NewOrderReadModel(
	orderID uuid.UUID,
	items []*ShopItemReadModel,
	accountEmail string,
	deliveryAddress string,
	deliveryTime time.Time,
) *OrderReadModel {
	return &OrderReadModel{
		ID: uuid.NewV4().
			String(),
		// we generate id ourself because auto generate mongo string id column with type _id is not an uuid
		OrderID:         orderID.String(),
		ShopItems:       items,
		AccountEmail:    accountEmail,
		DeliveryAddress: deliveryAddress,
		TotalPrice:      getShopItemsTotalPrice(items),
		DeliveredTime:   deliveryTime,
		CreatedAt:       time.Now(),
	}
}

// getShopItemsTotalPrice gets the total price of the shop items.
func getShopItemsTotalPrice(shopItems []*ShopItemReadModel) float64 {
	var totalPrice float64
	for _, item := range shopItems {
		totalPrice += item.Price * float64(item.Quantity)
	}

	return totalPrice
}
