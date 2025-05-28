package app

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
)

// CatalogsReadApplicationBuilder is a struct that contains the catalogs read application builder.
type CatalogsReadApplicationBuilder struct {
	contracts.ApplicationBuilder
}

// NewCatalogsReadApplicationBuilder is a constructor for the CatalogsReadApplicationBuilder.
func NewCatalogsReadApplicationBuilder() *CatalogsReadApplicationBuilder {
	builder := &CatalogsReadApplicationBuilder{fxapp.NewApplicationBuilder()}

	return builder
}

// Build is a method that builds the catalogs read application.
func (a *CatalogsReadApplicationBuilder) Build() *CatalogsReadApplication {
	return NewCatalogsReadApplication(
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
