package app

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
)

// CatalogReadApplicationBuilder is a struct that contains the catalogs read application builder.
type CatalogReadApplicationBuilder struct {
	contracts.ApplicationBuilder
}

// NewCatalogReadApplicationBuilder is a constructor for the CatalogReadApplicationBuilder.
func NewCatalogReadApplicationBuilder() *CatalogReadApplicationBuilder {
	builder := &CatalogReadApplicationBuilder{fxapp.NewApplicationBuilder()}

	return builder
}

// Build is a method that builds the catalogs read application.
func (a *CatalogReadApplicationBuilder) Build() *CatalogReadApplication {
	return NewCatalogReadApplication(
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
