package v1

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
)

// Ref: https://golangbot.com/inheritance/

type GetProducts struct {
	*utils.ListQuery
}

func NewGetProducts(query *utils.ListQuery) (*GetProducts, error) {
	return &GetProducts{ListQuery: query}, nil
}
