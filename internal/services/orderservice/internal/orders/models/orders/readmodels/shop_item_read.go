// Package readmodels contains the shop item read model.
package readmodels

// ShopItemReadModel is the read model for the shop item.
type ShopItemReadModel struct {
	Title       string  `json:"title,omitempty"       bson:"title,omitempty"`
	Description string  `json:"description,omitempty" bson:"description,omitempty"`
	Quantity    uint64  `json:"quantity,omitempty"    bson:"quantity,omitempty"`
	Price       float64 `json:"price,omitempty"       bson:"price,omitempty"`
}

// NewShopItemReadModel creates a new shop item read model.
func NewShopItemReadModel(
	title string,
	description string,
	quantity uint64,
	price float64,
) *ShopItemReadModel {
	return &ShopItemReadModel{
		Title:       title,
		Description: description,
		Quantity:    quantity,
		Price:       price,
	}
}
