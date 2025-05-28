package queries

import "github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

// Ref: https://golangbot.com/inheritance/

// GetProducts is a struct that contains the get products query.
type GetProducts struct {
	*utils.ListQuery
}

// NewGetProducts creates a new GetProducts.
func NewGetProducts(query *utils.ListQuery) *GetProducts {
	return &GetProducts{ListQuery: query}
}
