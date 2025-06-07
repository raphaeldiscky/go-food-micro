// Package integration contains the integration test fixture.
package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/bus"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"go.uber.org/zap"
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
	apptest "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/app/test"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/data/dbcontext"
	productsService "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/grpc/genproto"
)

// CatalogContext is an interface for integration tests that provides access to a product repository.
type CatalogContext interface {
	Products() contracts.ProductRepository
}

type catalogContextImpl struct {
	productRepo contracts.ProductRepository
}

func (c *catalogContextImpl) Products() contracts.ProductRepository {
	return c.productRepo
}

// CatalogUnitOfWork is an interface for integration tests that executes a function (using a CatalogContext) within a unit of work.
type CatalogUnitOfWork interface {
	Do(ctx context.Context, fn func(CatalogContext) error) error
}

type catalogUnitOfWorkImpl struct {
	productRepo contracts.ProductRepository
	db          *gorm.DB
	logger      logger.Logger
}

// handleTransactionRollback handles rolling back a transaction and logging any errors.
func (u *catalogUnitOfWorkImpl) handleTransactionRollback(tx *gorm.DB) {
	if err := tx.Rollback().Error; err != nil {
		u.logger.Error("failed to rollback transaction", zap.Error(err))
	}
}

// createTransactionRepository creates a new repository instance that uses the transaction.
func (u *catalogUnitOfWorkImpl) createTransactionRepository(
	tx *gorm.DB,
) (contracts.ProductRepository, error) {
	repo, ok := u.productRepo.(*repositories.PostgresProductRepository)
	if !ok {
		return nil, fmt.Errorf("failed to cast product repository to PostgresProductRepository")
	}

	if repo.Log == nil || repo.Tracer == nil {
		return nil, fmt.Errorf("logger or tracer is nil")
	}

	return repositories.NewPostgresProductRepository(repo.Log, tx, repo.Tracer), nil
}

// executeWithContext executes the given function with context handling and panic recovery.
func (u *catalogUnitOfWorkImpl) executeWithContext(
	ctx context.Context,
	tx *gorm.DB,
	fn func(CatalogContext) error,
	done chan<- error,
) {
	defer func() {
		if r := recover(); r != nil {
			u.handleTransactionRollback(tx)
			done <- fmt.Errorf("panic recovered: %v", r)
		}
	}()

	txRepo, err := u.createTransactionRepository(tx)
	if err != nil {
		done <- err

		return
	}

	catalogCtx := &catalogContextImpl{
		productRepo: txRepo,
	}

	if err := fn(catalogCtx); err != nil {
		u.handleTransactionRollback(tx)

		done <- err

		return
	}

	// Check if context is done before committing
	select {
	case <-ctx.Done():
		u.handleTransactionRollback(tx)
		done <- fmt.Errorf("context canceled: %w", ctx.Err())

		return
	default:
		// Context is not done, proceed with commit
		if err := tx.Commit().Error; err != nil {
			u.handleTransactionRollback(tx)
			done <- fmt.Errorf("failed to commit transaction: %w", err)

			return
		}
		done <- nil
	}
}

// Do executes the given function within a transaction, providing a CatalogContext that uses the transaction.
// It handles transaction management, context cancellation, and panic recovery.
// Returns an error if the transaction fails, context is canceled, or the function panics.
func (u *catalogUnitOfWorkImpl) Do(ctx context.Context, fn func(CatalogContext) error) error {
	tx := u.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	done := make(chan error, 1)
	go u.executeWithContext(ctx, tx, fn, done)

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		u.handleTransactionRollback(tx)

		return fmt.Errorf("context canceled: %w", ctx.Err())
	}
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
	result := apptest.NewCatalogWriteTestApp().Run(t)

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
			db:     result.Gorm,
			logger: result.Logger,
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

// NewCatalogUnitOfWork creates a new instance of CatalogUnitOfWork.
func NewCatalogUnitOfWork(result *apptest.CatalogWriteTestAppResult) CatalogUnitOfWork {
	// Create a no-op tracer for tests
	noopTracer := tracing.NewAppTracer("test")

	return &catalogUnitOfWorkImpl{
		productRepo: repositories.NewPostgresProductRepository(
			result.Logger,
			result.Gorm,
			noopTracer,
		),
		db:     result.Gorm,
		logger: result.Logger,
	}
}

// WaitForGrpcServerReady waits for the gRPC server to be ready.
func (i *CatalogWriteIntegrationTestSharedFixture) WaitForGrpcServerReady() error {
	if client, ok := i.ProductServiceClient.(interface{ WaitForAvailableConnection() error }); ok {
		return client.WaitForAvailableConnection()
	}

	return nil
}
