package test

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/test"

	"go.uber.org/fx/fxtest"
)

// CatalogsWriteTestApplicationBuilder is a struct that contains the catalogs write test application builder.
type CatalogsWriteTestApplicationBuilder struct {
	contracts.ApplicationBuilder
	tb fxtest.TB
}

// NewCatalogsWriteTestApplicationBuilder is a constructor for the CatalogsWriteTestApplicationBuilder.
func NewCatalogsWriteTestApplicationBuilder(tb fxtest.TB) *CatalogsWriteTestApplicationBuilder {
	return &CatalogsWriteTestApplicationBuilder{
		ApplicationBuilder: test.NewTestApplicationBuilder(tb),
		tb:                 tb,
	}
}

// Build is a method that builds the catalogs write test application.
func (a *CatalogsWriteTestApplicationBuilder) Build() *CatalogsWriteTestApplication {
	return NewCatalogsWriteTestApplication(
		a.tb,
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
