// Package test contains the test app.
package test

import (
	"context"
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/grpc"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/bus"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/testcontainer/gorm"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/testcontainer/rabbitmq"
	"github.com/stretchr/testify/require"

	fxcontracts "github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	config3 "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/config"
	contracts2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/migration/contracts"
	gormPostgres "github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm"
	config2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
	gorm2 "gorm.io/gorm"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/config"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/configurations/catalogs"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/data/dbcontext"
	productsService "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/grpc/genproto"
)

// CatalogWriteTestApp is a struct that contains the test app.
type CatalogWriteTestApp struct{}

// CatalogWriteTestAppResult is a struct that contains the test app result.
type CatalogWriteTestAppResult struct {
	Cfg                     *config.AppOptions
	Bus                     bus.RabbitmqBus
	Container               fxcontracts.Container
	Logger                  logger.Logger
	RabbitmqOptions         *config2.RabbitmqOptions
	EchoHTTPOptions         *config3.EchoHTTPOptions
	GormOptions             *gormPostgres.GormOptions
	Gorm                    *gorm2.DB
	ProductServiceClient    productsService.ProductsServiceClient
	GrpcClient              grpc.GrpcClient
	PostgresMigrationRunner contracts2.PostgresMigrationRunner
	CatalogsDBContext       *dbcontext.CatalogsGormDBContext
}

// NewCatalogWriteTestApp is a constructor for the CatalogWriteTestApp.
func NewCatalogWriteTestApp() *CatalogWriteTestApp {
	return &CatalogWriteTestApp{}
}

// Run is a method that runs the test app.
func (a *CatalogWriteTestApp) Run(t *testing.T) (result *CatalogWriteTestAppResult) {
	t.Helper()

	lifetimeCtx := context.Background()

	// ref: https://github.com/uber-go/fx/blob/master/app_test.go
	appBuilder := NewCatalogsWriteTestApplicationBuilder(t)
	appBuilder.ProvideModule(catalogs.NewCatalogsServiceModule())

	appBuilder.Decorate(
		rabbitmq.RabbitmqContainerOptionsDecorator(t, lifetimeCtx),
	)
	appBuilder.Decorate(gorm.GormContainerOptionsDecorator(t, lifetimeCtx))

	testApp := appBuilder.Build()

	err := testApp.ConfigureCatalogs()
	if err != nil {
		testApp.Logger().Fatalf("Error in ConfigureCatalogs, %s", err)
	}

	err = testApp.MapCatalogsEndpoints()
	if err != nil {
		testApp.Logger().Fatalf("Error in MapCatalogsEndpoints, %s", err)
	}

	testApp.ResolveFunc(
		func(cfg *config.AppOptions,
			bus bus.RabbitmqBus,
			logger logger.Logger,
			rabbitmqOptions *config2.RabbitmqOptions,
			gormOptions *gormPostgres.GormOptions,
			gorm *gorm2.DB,
			catalogsDBContext *dbcontext.CatalogsGormDBContext,
			echoOptions *config3.EchoHTTPOptions,
			grpcClient grpc.GrpcClient,
			postgresMigrationRunner contracts2.PostgresMigrationRunner,
		) {
			grpcConnection := grpcClient.GetGrpcConnection()

			result = &CatalogWriteTestAppResult{
				Bus:                     bus,
				Cfg:                     cfg,
				Container:               testApp,
				Logger:                  logger,
				RabbitmqOptions:         rabbitmqOptions,
				GormOptions:             gormOptions,
				Gorm:                    gorm,
				CatalogsDBContext:       catalogsDBContext,
				EchoHTTPOptions:         echoOptions,
				PostgresMigrationRunner: postgresMigrationRunner,
				ProductServiceClient: productsService.NewProductsServiceClient(
					grpcConnection,
				),
				GrpcClient: grpcClient,
			}
		},
	)
	// we need a longer timout for up and running our testcontainers
	duration := time.Second * 300

	// short timeout for handling start hooks and setup dependencies
	startCtx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	err = testApp.Start(startCtx)
	if err != nil {
		t.Fatalf("Error starting, err: %v", err)
	}

	// waiting for grpc endpoint becomes ready in the given timeout
	err = result.GrpcClient.WaitForAvailableConnection()
	require.NoError(t, err)

	t.Cleanup(func() {
		// short timeout for handling stop hooks
		stopCtx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		err = testApp.Stop(stopCtx)
		require.NoError(t, err)
	})

	return result
}
