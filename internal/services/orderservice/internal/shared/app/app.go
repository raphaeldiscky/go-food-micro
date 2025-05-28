// Package app contains the app.
package app

import "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/configurations/orders"

// App is a struct that contains the app.
type App struct{}

// NewApp creates a new App.
func NewApp() *App {
	return &App{}
}

// Run runs the app.
func (a *App) Run() {
	// configure dependencies
	appBuilder := NewOrdersApplicationBuilder()
	appBuilder.ProvideModule(orders.OrderServiceModule)

	app := appBuilder.Build()

	// configure application
	app.ConfigureOrders()

	app.MapOrdersEndpoints()

	app.Logger().Info("Starting orders_service application")
	app.Run()
}
