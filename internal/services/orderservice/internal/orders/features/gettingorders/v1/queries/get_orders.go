// Package queries contains the queries for the get orders.
package queries

import "github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

// Ref: https://golangbot.com/inheritance/

// GetOrders is the query for the get orders.
type GetOrders struct {
	*utils.ListQuery
}

// NewGetOrders creates a new get orders query.
func NewGetOrders(query *utils.ListQuery) *GetOrders {
	return &GetOrders{ListQuery: query}
}
