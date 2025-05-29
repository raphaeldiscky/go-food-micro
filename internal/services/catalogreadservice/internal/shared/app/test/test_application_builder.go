package test

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/test"
	"go.uber.org/fx/fxtest"
)

// CatalogReadTestApplicationBuilder is a struct that contains the catalogs read test application builder.
type CatalogReadTestApplicationBuilder struct {
	contracts.ApplicationBuilder
	tb fxtest.TB
}

// NewCatalogReadTestApplicationBuilder creates a new CatalogReadTestApplicationBuilder.
func NewCatalogReadTestApplicationBuilder(tb fxtest.TB) *CatalogReadTestApplicationBuilder {
	return &CatalogReadTestApplicationBuilder{
		ApplicationBuilder: test.NewTestApplicationBuilder(tb),
		tb:                 tb,
	}
}

// Build is a method that builds the catalogs read test application.
func (a *CatalogReadTestApplicationBuilder) Build() *CatalogReadTestApplication {
	return NewCatalogReadTestApplication(
		a.tb,
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
