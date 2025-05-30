package app

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
)

// OrdersApplicationBuilder is a struct that contains the orders application builder.
type OrdersApplicationBuilder struct {
	contracts.ApplicationBuilder
}

// NewOrdersApplicationBuilder is a constructor for the OrdersApplicationBuilder.
func NewOrdersApplicationBuilder() *OrdersApplicationBuilder {
	return &OrdersApplicationBuilder{fxapp.NewApplicationBuilder()}
}

// Build is a method that builds the orders application.
func (a *OrdersApplicationBuilder) Build() *OrderServiceApplication {
	return NewOrderServiceApplication(
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
