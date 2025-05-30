package app

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
)

// CatalogsWriteApplicationBuilder is a struct that contains the catalogs write application builder.
type CatalogsWriteApplicationBuilder struct {
	contracts.ApplicationBuilder
}

// NewCatalogsWriteApplicationBuilder is a constructor for the CatalogsWriteApplicationBuilder.
func NewCatalogsWriteApplicationBuilder() *CatalogsWriteApplicationBuilder {
	builder := &CatalogsWriteApplicationBuilder{fxapp.NewApplicationBuilder()}

	return builder
}

// Build is a method that builds the catalogs write application.
func (a *CatalogsWriteApplicationBuilder) Build() *CatalogsWriteApplication {
	return NewCatalogsWriteApplication(
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
