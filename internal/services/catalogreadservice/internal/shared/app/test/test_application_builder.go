package test

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/test"

	"go.uber.org/fx/fxtest"
)

// CatalogsReadTestApplicationBuilder is a struct that contains the catalogs read test application builder.
type CatalogsReadTestApplicationBuilder struct {
	contracts.ApplicationBuilder
	tb fxtest.TB
}

// NewCatalogsReadTestApplicationBuilder is a constructor for the CatalogsReadTestApplicationBuilder.
func NewCatalogsReadTestApplicationBuilder(tb fxtest.TB) *CatalogsReadTestApplicationBuilder {
	return &CatalogsReadTestApplicationBuilder{
		ApplicationBuilder: test.NewTestApplicationBuilder(tb),
		tb:                 tb,
	}
}

// Build is a method that builds the catalogs read test application.
func (a *CatalogsReadTestApplicationBuilder) Build() *CatalogsReadTestApplication {
	return NewCatalogsReadTestApplication(
		a.tb,
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
