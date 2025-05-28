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
type CatalogsReadTestApplication struct {
	*app.CatalogsReadApplication
	tb fxtest.TB
}

// NewCatalogsReadTestApplication is a constructor for the CatalogsReadTestApplication.
func NewCatalogsReadTestApplication(
	tb fxtest.TB,
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	logger logger.Logger,
	environment environment.Environment,
) *CatalogsReadTestApplication {
	testApp := test.NewTestApplication(
		tb,
		providers,
		decorates,
		options,
		logger,
		environment,
	)

	catalogApplication := &app.CatalogsReadApplication{
		CatalogsServiceConfigurator: catalogs.NewCatalogsServiceConfigurator(testApp),
	}

	return &CatalogsReadTestApplication{
		CatalogsReadApplication: catalogApplication,
		tb:                      tb,
	}
}
