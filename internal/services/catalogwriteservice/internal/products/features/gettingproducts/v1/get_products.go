package v1

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
)

// Ref: https://golangbot.com/inheritance/

// GetProducts is a struct that contains the get products query.
type GetProducts struct {
	*utils.ListQuery
}

// NewGetProducts is a constructor for the GetProducts.
func NewGetProducts(query *utils.ListQuery) (*GetProducts, error) {
	return &GetProducts{ListQuery: query}, nil
}
