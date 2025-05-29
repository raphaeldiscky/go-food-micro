// Package valueobject contains the shop item value object.
package valueobject

import (
	"fmt"
)

// ShopItem is the value object for the shop item.
type ShopItem struct {
	title       string
	description string
	quantity    uint64
	price       float64
}

// CreateNewShopItem creates a new shop item.
func CreateNewShopItem(title, description string, quantity uint64, price float64) *ShopItem {
	return &ShopItem{
		title:       title,
		description: description,
		quantity:    quantity,
		price:       price,
	}
}

// Title returns the title of the shop item.
func (s *ShopItem) Title() string {
	return s.title
}

// Description returns the description of the shop item.
func (s *ShopItem) Description() string {
	return s.description
}

// Quantity returns the quantity of the shop item.
func (s *ShopItem) Quantity() uint64 {
	return s.quantity
}

// Price returns the price of the shop item.
func (s *ShopItem) Price() float64 {
	return s.price
}

// String returns the string representation of the shop item.
func (s *ShopItem) String() string {
	return fmt.Sprintf("Title: {%s}, Description: {%s}, Quantity: {%v}, Price: {%v},",
		s.title,
		s.description,
		s.quantity,
		s.price,
	)
}
