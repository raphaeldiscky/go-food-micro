package test

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/test"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/app"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/configurations/catalogs"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// CatalogsWriteTestApplication is a struct that contains the catalogs write test application.
type CatalogsWriteTestApplication struct {
	*app.CatalogsWriteApplication
	tb fxtest.TB
}

// NewCatalogsWriteTestApplication is a constructor for the CatalogsWriteTestApplication.
func NewCatalogsWriteTestApplication(
	tb fxtest.TB,
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	logger logger.Logger,
	environment environment.Environment,
) *CatalogsWriteTestApplication {
	testApp := test.NewTestApplication(
		tb,
		providers,
		decorates,
		options,
		logger,
		environment,
	)

	catalogApplication := &app.CatalogsWriteApplication{
		CatalogsServiceConfigurator: catalogs.NewCatalogsServiceConfigurator(testApp),
	}

	return &CatalogsWriteTestApplication{
		CatalogsWriteApplication: catalogApplication,
		tb:                       tb,
	}
}
