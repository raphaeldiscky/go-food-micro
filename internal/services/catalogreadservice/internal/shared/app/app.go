// Package app contains the app.
package app

import "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/configurations/catalogs"

// CatalogReadApp is a struct that contains the app.
type CatalogReadApp struct{}

// NewCatalogReadApp creates a new CatalogReadApp.
func NewCatalogReadApp() *CatalogReadApp {
	return &CatalogReadApp{}
}

// Run runs the app.
func (a *CatalogReadApp) Run() {
	// configure dependencies
	appBuilder := NewCatalogReadApplicationBuilder()
	appBuilder.ProvideModule(catalogs.NewCatalogsServiceModule())

	app := appBuilder.Build()

	// configure application
	app.ConfigureCatalogs()

	app.MapCatalogsEndpoints()

	app.Logger().Info("Starting catalog_service application")
	app.Run()
}
