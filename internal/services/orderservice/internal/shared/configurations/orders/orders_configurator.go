// Package orders contains the orders configurator.
package orders

import (
	"fmt"
	"net/http"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"

	echo "github.com/labstack/echo/v4"
	echocontracts "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/contracts"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/config"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/configurations"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/configurations/orders/infrastructure"
)

// OrdersServiceConfigurator is the orders service configurator.
type OrdersServiceConfigurator struct {
	contracts.Application
	infrastructureConfigurator *infrastructure.InfrastructureConfigurator
	ordersModuleConfigurator   *configurations.OrdersModuleConfigurator
}

// NewOrdersServiceConfigurator creates a new orders service configurator.
func NewOrdersServiceConfigurator(
	app contracts.Application,
) *OrdersServiceConfigurator {
	infraConfigurator := infrastructure.NewInfrastructureConfigurator(app)
	ordersModuleConfigurator := configurations.NewOrdersModuleConfigurator(app)

	return &OrdersServiceConfigurator{
		Application:                app,
		infrastructureConfigurator: infraConfigurator,
		ordersModuleConfigurator:   ordersModuleConfigurator,
	}
}

// ConfigureOrders configures the orders.
func (ic *OrdersServiceConfigurator) ConfigureOrders() {
	ic.infrastructureConfigurator.ConfigInfrastructures()
	ic.ordersModuleConfigurator.ConfigureOrdersModule()
}

// MapOrdersEndpoints maps the orders endpoints.
func (ic *OrdersServiceConfigurator) MapOrdersEndpoints() {
	// Shared
	ic.ResolveFunc(
		func(ordersServer echocontracts.EchoHttpServer, cfg *config.Config) error {
			ordersServer.SetupDefaultMiddlewares()

			// config orders root endpoint
			ordersServer.RouteBuilder().
				RegisterRoutes(func(e *echo.Echo) {
					e.GET("", func(ec echo.Context) error {
						return ec.String(
							http.StatusOK,
							fmt.Sprintf(
								"%s is running...",
								cfg.AppOptions.GetMicroserviceNameUpper(),
							),
						)
					})
				})

			// config orders swagger
			ic.configSwagger(ordersServer.RouteBuilder())

			return nil
		},
	)

	// Modules
	// Orders Module endpoints
	ic.ordersModuleConfigurator.MapOrdersEndpoints()
}
