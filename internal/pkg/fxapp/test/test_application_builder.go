// Package test provides a set of functions for the test application builder.
package test

import (
	"go.uber.org/fx/fxtest"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
)

// TestApplicationBuilder is a struct that contains the test application builder.
type TestApplicationBuilder struct {
	contracts.ApplicationBuilder
	TB fxtest.TB
}

// NewTestApplicationBuilder creates a new test application builder.
func NewTestApplicationBuilder(tb fxtest.TB) *TestApplicationBuilder {
	return &TestApplicationBuilder{
		TB:                 tb,
		ApplicationBuilder: fxapp.NewApplicationBuilder(environment.Test),
	}
}

// Build builds the test application.
func (a *TestApplicationBuilder) Build() contracts.Application {
	app := NewTestApplication(
		a.TB,
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		environment.Test,
	)

	return app
}
