package test

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/test"
	"go.uber.org/fx/fxtest"
)

// OrdersTestApplicationBuilder is a struct that contains the orders test application builder.
type OrdersTestApplicationBuilder struct {
	contracts.ApplicationBuilder
	tb fxtest.TB
}

// NewOrdersTestApplicationBuilder is a constructor for the OrdersTestApplicationBuilder.
func NewOrdersTestApplicationBuilder(tb fxtest.TB) *OrdersTestApplicationBuilder {
	return &OrdersTestApplicationBuilder{
		ApplicationBuilder: test.NewTestApplicationBuilder(tb),
		tb:                 tb,
	}
}

// Build is a method that builds the orders test application.
func (a *OrdersTestApplicationBuilder) Build() *OrdersTestApplication {
	return NewOrdersTestApplication(
		a.tb,
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
