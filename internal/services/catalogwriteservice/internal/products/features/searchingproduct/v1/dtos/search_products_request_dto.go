package dtos

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
)

// SearchProductsRequestDto is a struct that contains the search products request dto.
type SearchProductsRequestDto struct {
	SearchText       string `query:"search" json:"search"`
	*utils.ListQuery `                      json:"listQuery"`
}
