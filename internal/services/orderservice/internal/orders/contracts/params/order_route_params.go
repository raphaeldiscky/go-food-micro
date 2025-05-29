// Package params contains the order route parameters.
package params

import (
	"github.com/go-playground/validator"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"go.uber.org/fx"

	echo "github.com/labstack/echo/v4"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/contracts"
)

// OrderRouteParams is the parameters for the order route.
type OrderRouteParams struct {
	fx.In

	OrdersMetrics *contracts.OrdersMetrics
	Logger        logger.Logger
	OrdersGroup   *echo.Group `name:"order-echo-group"`
	Validator     *validator.Validate
}
