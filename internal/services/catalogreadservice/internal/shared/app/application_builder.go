package app

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
)

type CatalogsReadApplicationBuilder struct {
	contracts.ApplicationBuilder
}

func NewCatalogsReadApplicationBuilder() *CatalogsReadApplicationBuilder {
	builder := &CatalogsReadApplicationBuilder{fxapp.NewApplicationBuilder()}

	return builder
}

func (a *CatalogsReadApplicationBuilder) Build() *CatalogsReadApplication {
	return NewCatalogsReadApplication(
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
