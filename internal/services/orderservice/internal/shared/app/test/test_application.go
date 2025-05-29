package test

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/test"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/app"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/configurations/orders"
)

// OrderServiceTestApplication is a struct that contains the orders test application.
type OrderServiceTestApplication struct {
	*app.OrderServiceApplication
	tb fxtest.TB
}

// NewOrderServiceTestApplication creates a new OrderServiceTestApplication.
func NewOrderServiceTestApplication(
	tb fxtest.TB,
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	log logger.Logger,
	env environment.Environment,
) *OrderServiceTestApplication {
	testApp := test.NewTestApplication(
		tb,
		providers,
		decorates,
		options,
		log,
		env,
	)

	orderApplication := &app.OrderServiceApplication{
		OrderServiceConfigurator: orders.NewOrderServiceConfigurator(testApp),
	}

	return &OrderServiceTestApplication{
		OrderServiceApplication: orderApplication,
		tb:                      tb,
	}
}
