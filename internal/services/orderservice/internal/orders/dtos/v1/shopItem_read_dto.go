// Package dtosv1 contains the shop item read dto.
package dtosv1

// ShopItemReadDto is the read dto for the shop item.
type ShopItemReadDto struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Quantity    uint64  `json:"quantity"`
	Price       float64 `json:"price"`
}
