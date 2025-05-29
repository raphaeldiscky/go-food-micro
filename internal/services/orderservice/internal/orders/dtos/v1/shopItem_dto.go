// Package dtosv1 contains the shop item dto.
package dtosv1

// ShopItemDto is the dto for the shop item.
type ShopItemDto struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Quantity    uint64  `json:"quantity"`
	Price       float64 `json:"price"`
}
