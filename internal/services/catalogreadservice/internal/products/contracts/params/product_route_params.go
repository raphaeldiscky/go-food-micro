package params

import (
	"github.com/go-playground/validator"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"go.uber.org/fx"

	echo "github.com/labstack/echo/v4"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/contracts"
)

type ProductRouteParams struct {
	fx.In

	CatalogsMetrics *contracts.CatalogsMetrics
	Logger          logger.Logger
	ProductsGroup   *echo.Group `name:"product-echo-group"`
	Validator       *validator.Validate
}
