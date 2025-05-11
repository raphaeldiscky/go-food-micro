package dtos

import dtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/dtos/v1"

type GetOrderByIdResponseDto struct {
	Order *dtosV1.OrderReadDto `json:"order"`
}
