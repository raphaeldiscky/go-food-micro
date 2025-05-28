// Package app contains the app.
package app

import "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/configurations/catalogs"

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	// configure dependencies
	appBuilder := NewCatalogsReadApplicationBuilder()
	appBuilder.ProvideModule(catalogs.NewCatalogsServiceModule())

	app := appBuilder.Build()

	// configure application
	app.ConfigureCatalogs()

	app.MapCatalogsEndpoints()

	app.Logger().Info("Starting catalog_service application")
	app.Run()
}
