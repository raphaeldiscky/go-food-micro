package app

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/configurations/orders"
)

// OrdersApplication is a struct that contains the orders application.
type OrdersApplication struct {
	*orders.OrdersServiceConfigurator
}

// NewOrdersApplication creates a new OrdersApplication.
func NewOrdersApplication(
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	logger logger.Logger,
	environment environment.Environment,
) *OrdersApplication {
	app := fxapp.NewApplication(providers, decorates, options, logger, environment)

	return &OrdersApplication{
		OrdersServiceConfigurator: orders.NewOrdersServiceConfigurator(app),
	}
}
