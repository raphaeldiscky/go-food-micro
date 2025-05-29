// Package dtos contains the get order by id response dto.
package dtos

import dtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/dtos/v1"

// GetOrderByIDResponseDto is the response dto for the get order by id endpoint.
type GetOrderByIDResponseDto struct {
	Order *dtosV1.OrderReadDto `json:"order"`
}
