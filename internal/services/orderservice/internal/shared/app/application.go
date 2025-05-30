package app

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/configurations/orders"
)

// OrderServiceApplication is a struct that contains the orders application.
type OrderServiceApplication struct {
	*orders.OrderServiceConfigurator
}

// NewOrderServiceApplication creates a new OrderServiceApplication.
func NewOrderServiceApplication(
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	log logger.Logger,
	env environment.Environment,
) *OrderServiceApplication {
	app := fxapp.NewApplication(providers, decorates, options, log, env)

	return &OrderServiceApplication{
		OrderServiceConfigurator: orders.NewOrderServiceConfigurator(app),
	}
}
