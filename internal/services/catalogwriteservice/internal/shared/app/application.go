package app

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/configurations/catalogs"

	"go.uber.org/fx"
)

// CatalogsWriteApplication is a struct that contains the catalogs write application.
type CatalogsWriteApplication struct {
	*catalogs.CatalogsServiceConfigurator
}

// NewCatalogsWriteApplication is a constructor for the CatalogsWriteApplication.
func NewCatalogsWriteApplication(
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	logger logger.Logger,
	environment environment.Environment,
) *CatalogsWriteApplication {
	app := fxapp.NewApplication(providers, decorates, options, logger, environment)
	return &CatalogsWriteApplication{
		CatalogsServiceConfigurator: catalogs.NewCatalogsServiceConfigurator(app),
	}
}
