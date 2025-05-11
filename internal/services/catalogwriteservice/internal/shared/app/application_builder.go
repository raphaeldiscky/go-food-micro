package app

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
)

type CatalogsWriteApplicationBuilder struct {
	contracts.ApplicationBuilder
}

func NewCatalogsWriteApplicationBuilder() *CatalogsWriteApplicationBuilder {
	builder := &CatalogsWriteApplicationBuilder{fxapp.NewApplicationBuilder()}

	return builder
}

func (a *CatalogsWriteApplicationBuilder) Build() *CatalogsWriteApplication {
	return NewCatalogsWriteApplication(
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
