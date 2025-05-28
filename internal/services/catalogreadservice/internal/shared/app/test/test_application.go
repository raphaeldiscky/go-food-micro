package test

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/test"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/app"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/configurations/catalogs"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// CatalogsReadTestApplication is a struct that contains the catalogs read test application.
type CatalogReadTestApplication struct {
	*app.CatalogReadApplication
	tb fxtest.TB
}

// NewCatalogsReadTestApplication is a constructor for the CatalogsReadTestApplication.
func NewCatalogReadTestApplication(
	tb fxtest.TB,
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	log logger.Logger,
	env environment.Environment,
) *CatalogReadTestApplication {
	testApp := test.NewTestApplication(
		tb,
		providers,
		decorates,
		options,
		log,
		env,
	)

	catalogApplication := &app.CatalogReadApplication{
		CatalogReadServiceConfigurator: catalogs.NewCatalogReadServiceConfigurator(testApp),
	}

	return &CatalogReadTestApplication{
		CatalogReadApplication: catalogApplication,
		tb:                     tb,
	}
}
