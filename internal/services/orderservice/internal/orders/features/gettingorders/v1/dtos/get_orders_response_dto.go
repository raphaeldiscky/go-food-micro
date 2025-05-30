// Package dtos contains the dtos for the get orders.
package dtos

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	dtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/dtos/v1"
)

// GetOrdersResponseDto is the response dto for the get orders.
type GetOrdersResponseDto struct {
	Orders *utils.ListResult[*dtosV1.OrderReadDto]
}
