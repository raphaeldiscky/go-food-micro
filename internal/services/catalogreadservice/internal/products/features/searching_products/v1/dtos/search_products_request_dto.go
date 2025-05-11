package dtos

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
)

type SearchProductsRequestDto struct {
	SearchText       string `query:"search" json:"search"`
	*utils.ListQuery `                      json:"listQuery"`
}
