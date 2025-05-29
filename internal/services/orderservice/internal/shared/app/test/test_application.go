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

// OrdersTestApplication is a struct that contains the orders test application.
type OrdersTestApplication struct {
	*app.OrdersApplication
	tb fxtest.TB
}

// NewOrdersTestApplication creates a new OrdersTestApplication.
func NewOrdersTestApplication(
	tb fxtest.TB,
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	logger logger.Logger,
	env environment.Environment,
) *OrdersTestApplication {
	testApp := test.NewTestApplication(
		tb,
		providers,
		decorates,
		options,
		logger,
		env,
	)

	orderApplication := &app.OrdersApplication{
		OrdersServiceConfigurator: orders.NewOrdersServiceConfigurator(testApp),
	}

	return &OrdersTestApplication{
		OrdersApplication: orderApplication,
		tb:                tb,
	}
}
