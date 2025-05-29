// Package dtos contains the get products request dto.
package dtos

import "github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

// GetProductsRequestDto is a struct that contains the get products request dto.
type GetProductsRequestDto struct {
	*utils.ListQuery
}
