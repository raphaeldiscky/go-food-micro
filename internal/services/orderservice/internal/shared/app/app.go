// Package app contains the app.
package app

import "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/configurations/orders"

// OrderApp is a struct that contains the app.
type OrderApp struct{}

// NewOrderApp creates a new OrderApp.
func NewOrderApp() *OrderApp {
	return &OrderApp{}
}

// Run runs the app.
func (a *OrderApp) Run() {
	// configure dependencies
	appBuilder := NewOrdersApplicationBuilder()
	appBuilder.ProvideModule(orders.OrderServiceModule())

	app := appBuilder.Build()

	// configure application
	app.ConfigureOrders()

	app.MapOrdersEndpoints()

	app.Logger().Info("Starting orders_service application")
	app.Run()
}
