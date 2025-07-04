// Package app contains the app.
package app

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/configurations/catalogs"
)

// CatalogWriteApp is a struct that contains the app.
type CatalogWriteApp struct{}

// NewCatalogWriteApp is a constructor for the CatalogWriteApp.
func NewCatalogWriteApp() *CatalogWriteApp {
	return &CatalogWriteApp{}
}

// Run is a method that runs the app.
func (a *CatalogWriteApp) Run() {
	// configure dependencies
	appBuilder := NewCatalogsWriteApplicationBuilder()
	appBuilder.ProvideModule(catalogs.NewCatalogsServiceModule())

	app := appBuilder.Build()

	// configure application
	err := app.ConfigureCatalogs()
	if err != nil {
		app.Logger().Fatalf("Error in ConfigureCatalogs", err)
	}

	err = app.MapCatalogsEndpoints()
	if err != nil {
		app.Logger().Fatalf("Error in MapCatalogsEndpoints", err)
	}

	app.Logger().Info("Starting catalog_service application")
	app.ResolveFunc(func(tracer tracing.AppTracer) {
		_, span := tracer.Start(context.Background(), "Application started")
		span.End()
	})

	app.Run()
}
