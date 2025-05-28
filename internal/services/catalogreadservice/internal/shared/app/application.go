package app

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/configurations/catalogs"
)

// CatalogReadApplication is a struct that contains the catalog read application.
type CatalogReadApplication struct {
	*catalogs.CatalogReadServiceConfigurator
}

// NewCatalogReadApplication creates a new CatalogReadApplication.
func NewCatalogReadApplication(
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	log logger.Logger,
	env environment.Environment,
) *CatalogReadApplication {
	app := fxapp.NewApplication(providers, decorates, options, log, env)

	return &CatalogReadApplication{
		CatalogReadServiceConfigurator: catalogs.NewCatalogReadServiceConfigurator(app),
	}
}
