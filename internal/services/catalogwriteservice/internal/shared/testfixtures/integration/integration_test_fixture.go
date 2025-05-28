// Package integration contains the integration test fixture.
package integration

import (
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/bus"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"gorm.io/gorm"

	_ "github.com/lib/pq" // postgres driver

	gofakeit "github.com/brianvoe/gofakeit/v6"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	fxcontracts "github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	config2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
	uuid "github.com/satori/go.uuid"
	dbcleaner "gopkg.in/khaiql/dbcleaner.v2"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/config"
	datamodel "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/data/datamodels"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/app/test"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/data/dbcontext"
	productsService "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/grpc/genproto"
)

// IntegrationTestSharedFixture is a struct that contains the integration test shared fixture.
type IntegrationTestSharedFixture struct {
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
}

// NewIntegrationTestSharedFixture is a constructor for the IntegrationTestSharedFixture.
func NewIntegrationTestSharedFixture(
	t *testing.T,
) *IntegrationTestSharedFixture {
	t.Helper()
	result := test.NewCatalogWriteTestApp().Run(t)

	// https://github.com/michaelklishin/rabbit-hole
	rmqc, err := rabbithole.NewClient(
		result.RabbitmqOptions.RabbitmqHostOptions.HttpEndPoint(),
		result.RabbitmqOptions.RabbitmqHostOptions.UserName,
		result.RabbitmqOptions.RabbitmqHostOptions.Password)
	if err != nil {
		result.Logger.Error(
			errors.WrapIf(err, "error in creating rabbithole client"),
		)
	}

	shared := &IntegrationTestSharedFixture{
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
	}

	return shared
}

// SetupTest is a method that sets up the test.
func (i *IntegrationTestSharedFixture) SetupTest() {
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
func (i *IntegrationTestSharedFixture) TearDownTest() {
	i.Log.Info("TearDownTest started")

	// cleanup test containers with their hooks
	if err := i.cleanupRabbitmqData(); err != nil {
		i.Log.Error(errors.WrapIf(err, "error in cleanup rabbitmq data"))
	}

	if err := i.cleanupPostgresData(); err != nil {
		i.Log.Error(errors.WrapIf(err, "error in cleanup postgres data"))
	}
}

func (i *IntegrationTestSharedFixture) cleanupRabbitmqData() error {
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

func (i *IntegrationTestSharedFixture) cleanupPostgresData() error {
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
