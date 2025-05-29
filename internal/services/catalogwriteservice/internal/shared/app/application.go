package app

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/configurations/catalogs"
)

// CatalogsWriteApplication is a struct that contains the catalogs write application.
type CatalogsWriteApplication struct {
	*catalogs.CatalogWriteServiceConfigurator
}

// NewCatalogsWriteApplication creates a new CatalogsWriteApplication.
func NewCatalogsWriteApplication(
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	log logger.Logger,
	env environment.Environment,
) *CatalogsWriteApplication {
	app := fxapp.NewApplication(providers, decorates, options, log, env)

	return &CatalogsWriteApplication{
		CatalogWriteServiceConfigurator: catalogs.NewCatalogWriteServiceConfigurator(app),
	}
}
