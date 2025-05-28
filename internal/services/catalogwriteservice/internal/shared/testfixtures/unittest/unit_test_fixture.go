// Package unittest contains the unit test fixture.
package unittest

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/glebarez/sqlite"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/mocks"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/external/gromlog"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/helpers/gormextensions"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	defaultLogger "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/config"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/configurations/mappings"
	datamodel "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/data/datamodels"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/data/dbcontext"
)

// UnitTestSharedFixture is a struct that contains the shared fixture for the unit tests.
type UnitTestSharedFixture struct {
	Cfg *config.AppOptions
	Log logger.Logger
	suite.Suite
	Products         []*datamodel.ProductDataModel
	Bus              *mocks.Bus
	Tracer           trace.Tracer
	CatalogDBContext *dbcontext.CatalogsGormDBContext
	Ctx              context.Context
	dbFilePath       string
	dbFileName       string
}

// NewUnitTestSharedFixture is a constructor for the UnitTestSharedFixture.
func NewUnitTestSharedFixture(_ *testing.T) *UnitTestSharedFixture {
	// we could use EmptyLogger if we don't want to log anything
	log := defaultLogger.GetLogger()
	cfg := &config.AppOptions{}

	// empty tracer, just for testing
	nopetracer := trace.NewNoopTracerProvider()
	testTracer := nopetracer.Tracer("test_tracer")

	unit := &UnitTestSharedFixture{
		Cfg:        cfg,
		Log:        log,
		Tracer:     testTracer,
		dbFileName: "sqlite.db",
	}

	return unit
}

// BeginTx is a method that begins a transaction.
func (c *UnitTestSharedFixture) BeginTx() {
	c.Log.Info("starting transaction")
	// seems when we `Begin` a transaction on gorm.DB (with SQLLite in-memory) our previous gormDB before transaction will remove and the new gormDB with tx will go on the memory
	tx := c.CatalogDBContext.DB().Begin()
	gormContext := gormextensions.SetTxToContext(c.Ctx, tx)
	c.Ctx = gormContext
}

// CommitTx is a method that commits the transaction.
func (c *UnitTestSharedFixture) CommitTx() {
	tx := gormextensions.GetTxFromContextIfExists(c.Ctx)
	if tx != nil {
		c.Log.Info("committing transaction")
		tx.Commit()
	}
}

// SetupSuite is a hook that is called before all tests in the suite have run.
func (c *UnitTestSharedFixture) SetupSuite() {
	// this fix root working directory problem in our test environment inner our fixture
	environment.FixProjectRootWorkingDirectoryPath()
	projectRootDir := environment.GetProjectRootWorkingDirectory()

	c.dbFilePath = filepath.Join(projectRootDir, c.dbFileName)
}

// TearDownSuite is a hook that is called after all tests in the suite have run.
func (c *UnitTestSharedFixture) TearDownSuite() {
}

// SetupTest is a hook that is called before each test.
func (c *UnitTestSharedFixture) SetupTest() {
	ctx := context.Background()
	c.Ctx = ctx

	c.setupBus()

	c.setupDB()

	err := mappings.ConfigureProductsMappings()
	c.Require().NoError(err)
}

// TearDownTest is a hook that is called after each test.
func (c *UnitTestSharedFixture) TearDownTest() {
	err := c.cleanupDB()
	c.Require().NoError(err)

	mapper.ClearMappings()
}

// setupBus is a method that sets up the bus.
func (c *UnitTestSharedFixture) setupBus() {
	bus := &mocks.Bus{}
	bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	c.Bus = bus
}

// setupDB is a method that sets up the database.
func (c *UnitTestSharedFixture) setupDB() {
	dbContext := c.createSQLLiteDBContext()
	c.CatalogDBContext = dbContext

	c.initDB(dbContext)
}

// createSQLLiteDBContext is a method that creates the SQLLite database context.
func (c *UnitTestSharedFixture) createSQLLiteDBContext() *dbcontext.CatalogsGormDBContext {
	// https://gorm.io/docs/connecting_to_the_database.html#SQLite
	// https://github.com/glebarez/sqlite
	// https://www.connectionstrings.com/sqlite/
	gormSQLLiteDB, err := gorm.Open(
		sqlite.Open(c.dbFilePath),
		&gorm.Config{
			Logger: gromlog.NewGormCustomLogger(defaultLogger.GetLogger()),
		})
	c.Require().NoError(err)

	dbContext := dbcontext.NewCatalogsDBContext(gormSQLLiteDB)

	return dbContext
}

// initDB is a method that initializes the database.
func (c *UnitTestSharedFixture) initDB(dbContext *dbcontext.CatalogsGormDBContext) {
	// migrations for our database
	err := migrateGorm(dbContext)
	c.Require().NoError(err)

	// seed data for our tests
	items, err := seedDataManually(dbContext)
	c.Require().NoError(err)

	c.Products = items
}

// cleanupDB is a method that cleans up the database.
func (c *UnitTestSharedFixture) cleanupDB() error {
	sqldb, err := c.CatalogDBContext.DB().DB()
	if err != nil {
		return err
	}
	err = sqldb.Close()
	if err != nil {
		return err
	}

	// removing sql-lite file
	err = os.Remove(c.dbFilePath)
	if err != nil {
		return err
	}

	return err
}

// migrateGorm is a method that migrates the Gorm database.
func migrateGorm(dbContext *dbcontext.CatalogsGormDBContext) error {
	err := dbContext.DB().AutoMigrate(&datamodel.ProductDataModel{})
	if err != nil {
		return err
	}

	return nil
}

// seedDataManually is a method that seeds the database with data.
func seedDataManually(
	dbContext *dbcontext.CatalogsGormDBContext,
) ([]*datamodel.ProductDataModel, error) {
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

	// seed data
	err := dbContext.DB().CreateInBatches(products, len(products)).Error
	if err != nil {
		return nil, errors.Wrap(err, "error in seed database")
	}

	return products, nil
}
