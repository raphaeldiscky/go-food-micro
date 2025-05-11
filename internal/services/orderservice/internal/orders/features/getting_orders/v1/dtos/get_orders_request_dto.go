package dtos

import "github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

type GetOrdersRequestDto struct {
	*utils.ListQuery
}
