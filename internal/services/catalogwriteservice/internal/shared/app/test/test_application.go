package test

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/test"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/app"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/configurations/catalogs"
)

// CatalogsWriteTestApplication is a struct that contains the catalogs write test application.
type CatalogsWriteTestApplication struct {
	*app.CatalogsWriteApplication
	tb fxtest.TB
}

// NewCatalogsWriteTestApplication creates a new CatalogsWriteTestApplication.
func NewCatalogsWriteTestApplication(
	tb fxtest.TB,
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	log logger.Logger,
	env environment.Environment,
) *CatalogsWriteTestApplication {
	testApp := test.NewTestApplication(
		tb,
		providers,
		decorates,
		options,
		log,
		env,
	)

	catalogApplication := &app.CatalogsWriteApplication{
		CatalogWriteServiceConfigurator: catalogs.NewCatalogWriteServiceConfigurator(testApp),
	}

	return &CatalogsWriteTestApplication{
		CatalogsWriteApplication: catalogApplication,
		tb:                       tb,
	}
}
