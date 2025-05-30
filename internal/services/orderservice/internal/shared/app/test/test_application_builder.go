package test

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/test"
	"go.uber.org/fx/fxtest"
)

// OrderServiceTestApplicationBuilder is a struct that contains the orders test application builder.
type OrderServiceTestApplicationBuilder struct {
	contracts.ApplicationBuilder
	tb fxtest.TB
}

// NewOrderServiceTestApplicationBuilder is a constructor for the OrderServiceTestApplicationBuilder.
func NewOrderServiceTestApplicationBuilder(tb fxtest.TB) *OrderServiceTestApplicationBuilder {
	return &OrderServiceTestApplicationBuilder{
		ApplicationBuilder: test.NewTestApplicationBuilder(tb),
		tb:                 tb,
	}
}

// Build is a method that builds the orders test application.
func (a *OrderServiceTestApplicationBuilder) Build() *OrderServiceTestApplication {
	return NewOrderServiceTestApplication(
		a.tb,
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
