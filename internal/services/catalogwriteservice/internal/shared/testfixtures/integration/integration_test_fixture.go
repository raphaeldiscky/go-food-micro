// Package integration contains the integration test fixture.
package integration

import (
	"context"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/bus"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"gorm.io/gorm"

	_ "github.com/lib/pq" // postgres driver

	gofakeit "github.com/brianvoe/gofakeit/v6"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	fxcontracts "github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	config2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
	uuid "github.com/satori/go.uuid"
	dbcleaner "gopkg.in/khaiql/dbcleaner.v2"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/config"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/contracts"
	datamodel "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/data/datamodels"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/data/repositories"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/app/test"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/data/dbcontext"
	productsService "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/grpc/genproto"
)

// CatalogContext provides access to repositories in a unit of work
// Only Products() is implemented for now
// You can expand this as needed for more repositories
// Place this near the top of the file

type CatalogContext interface {
	Products() contracts.ProductRepository
}

type catalogContextImpl struct {
	productRepo contracts.ProductRepository
}

func (c *catalogContextImpl) Products() contracts.ProductRepository {
	return c.productRepo
}

// CatalogUnitOfWork provides a Do method for transactional work
// Only a simple implementation is provided for integration tests

type CatalogUnitOfWork interface {
	Do(ctx context.Context, fn func(CatalogContext) error) error
}

type catalogUnitOfWorkImpl struct {
	productRepo contracts.ProductRepository
}

func (u *catalogUnitOfWorkImpl) Do(ctx context.Context, fn func(CatalogContext) error) error {
	// In a real implementation, you would start a DB transaction here
	// For integration tests, just call the function with the context
	return fn(&catalogContextImpl{productRepo: u.productRepo})
}

// CatalogWriteIntegrationTestSharedFixture is a struct that contains the integration test shared fixture.
type CatalogWriteIntegrationTestSharedFixture struct {
	Cfg                  *config.AppOptions
	Log                  logger.Logger
	Bus                  bus.Bus
	CatalogsDBContext    *dbcontext.CatalogsGormDBContext
	Container            fxcontracts.Container
	DbCleaner            dbcleaner.DbCleaner
	RabbitmqCleaner      *rabbithole.Client
	rabbitmqOptions      *config2.RabbitmqOptions
	Gorm                 *gorm.DB
	BaseAddress          string
	Items                []*datamodel.ProductDataModel
	ProductServiceClient productsService.ProductsServiceClient
	ProductRepository    contracts.ProductRepository
	tracer               tracing.AppTracer
	CatalogUnitOfWorks   CatalogUnitOfWork
}

// NewCatalogWriteIntegrationTestSharedFixture is a constructor for the CatalogWriteIntegrationTestSharedFixture.
func NewCatalogWriteIntegrationTestSharedFixture(
	t *testing.T,
) *CatalogWriteIntegrationTestSharedFixture {
	t.Helper()
	result := test.NewCatalogWriteTestApp().Run(t)

	// https://github.com/michaelklishin/rabbit-hole
	rmqc, err := rabbithole.NewClient(
		result.RabbitmqOptions.RabbitmqHostOptions.HTTPEndPoint(),
		result.RabbitmqOptions.RabbitmqHostOptions.UserName,
		result.RabbitmqOptions.RabbitmqHostOptions.Password)
	if err != nil {
		result.Logger.Error(
			errors.WrapIf(err, "error in creating rabbithole client"),
		)
	}

	// Create a no-op tracer for tests
	noopTracer := tracing.NewAppTracer("test")

	shared := &CatalogWriteIntegrationTestSharedFixture{
		Log:                  result.Logger,
		Container:            result.Container,
		Cfg:                  result.Cfg,
		RabbitmqCleaner:      rmqc,
		CatalogsDBContext:    result.CatalogsDBContext,
		Bus:                  result.Bus,
		rabbitmqOptions:      result.RabbitmqOptions,
		Gorm:                 result.Gorm,
		BaseAddress:          result.EchoHTTPOptions.BasePathAddress(),
		ProductServiceClient: result.ProductServiceClient,
		tracer:               noopTracer,
		ProductRepository: repositories.NewPostgresProductRepository(
			result.Logger,
			result.Gorm,
			noopTracer,
		),
		CatalogUnitOfWorks: &catalogUnitOfWorkImpl{
			productRepo: repositories.NewPostgresProductRepository(
				result.Logger,
				result.Gorm,
				noopTracer,
			),
		},
	}

	return shared
}

// SetupTest is a method that sets up the test.
func (i *CatalogWriteIntegrationTestSharedFixture) SetupTest() {
	i.Log.Info("SetupTest started")

	// migration will do in app configuration
	// seed data for our tests - app seed doesn't run in test environment
	res, err := seedDataManually(i.Gorm)
	if err != nil {
		i.Log.Error(errors.WrapIf(err, "error in seeding data in postgres"))
	}

	i.Items = res
}

// TearDownTest is a method that tears down the test.
func (i *CatalogWriteIntegrationTestSharedFixture) TearDownTest() {
	i.Log.Info("TearDownTest started")

	// cleanup test containers with their hooks
	if err := i.cleanupRabbitmqData(); err != nil {
		i.Log.Error(errors.WrapIf(err, "error in cleanup rabbitmq data"))
	}

	if err := i.cleanupPostgresData(); err != nil {
		i.Log.Error(errors.WrapIf(err, "error in cleanup postgres data"))
	}
}

func (i *CatalogWriteIntegrationTestSharedFixture) cleanupRabbitmqData() error {
	// https://github.com/michaelklishin/rabbit-hole
	// Get all queues
	queues, err := i.RabbitmqCleaner.ListQueuesIn(
		i.rabbitmqOptions.RabbitmqHostOptions.VirtualHost,
	)
	if err != nil {
		return err
	}
	// clear each queue
	var lastErr error
	for idx := range queues {
		_, err := i.RabbitmqCleaner.PurgeQueue(
			i.rabbitmqOptions.RabbitmqHostOptions.VirtualHost,
			queues[idx].Name,
		)
		if err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func (i *CatalogWriteIntegrationTestSharedFixture) cleanupPostgresData() error {
	tables := []string{"products"}
	// Iterate over the tables and delete all records
	for _, table := range tables {
		err := i.Gorm.Exec("DELETE FROM " + table).Error

		return err
	}

	return nil
}

func seedDataManually(gormDB *gorm.DB) ([]*datamodel.ProductDataModel, error) {
	products := []*datamodel.ProductDataModel{
		{
			ID:          uuid.NewV4(),
			Name:        gofakeit.Name(),
			CreatedAt:   time.Now(),
			Description: gofakeit.AdjectiveDescriptive(),
			Price:       gofakeit.Price(100, 1000),
		},
		{
			ID:          uuid.NewV4(),
			Name:        gofakeit.Name(),
			CreatedAt:   time.Now(),
			Description: gofakeit.AdjectiveDescriptive(),
			Price:       gofakeit.Price(100, 1000),
		},
	}

	err := gormDB.CreateInBatches(products, len(products)).Error
	if err != nil {
		return nil, errors.Wrap(err, "error in seed database")
	}

	return products, nil
}
