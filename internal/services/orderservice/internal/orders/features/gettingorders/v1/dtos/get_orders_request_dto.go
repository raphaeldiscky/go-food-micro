// Package dtos contains the dtos for the get orders.
package dtos

import "github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

// GetOrdersRequestDto is the request dto for the get orders.
type GetOrdersRequestDto struct {
	*utils.ListQuery
}
