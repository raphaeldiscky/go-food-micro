//go:build unit
// +build unit

package dbcontext

import (
	"context"
	"os"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/external/fxlog"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/zap"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/gormdbcontext"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/scopes"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	gormPostgres "github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/configurations/mappings"
	datamodel "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/data/datamodels"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/models"
)

// Define the suite
type DBContextTestSuite struct {
	suite.Suite
	items      []*datamodel.ProductDataModel
	dbContext  *CatalogsGormDBContext
	app        *fxtest.App
	dbFilePath string
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestDBContextTestSuite(t *testing.T) {
	suite.Run(t, new(DBContextTestSuite))
}

// TestFindProductByID tests the find product by id.
func (s *DBContextTestSuite) TestFindProductByID() {
	s.Require().NotNil(s.dbContext)

	id := s.items[0].ID

	p, err := gormdbcontext.FindModelByID[*datamodel.ProductDataModel, *models.Product](
		context.Background(),
		s.dbContext,
		id,
	)
	s.Require().NoError(err)
	s.Require().NotNil(p)

	s.Assert().Equal(p.ID, id)
}

// TestExistsProductByID tests the exists product by id.
func (s *DBContextTestSuite) TestExistsProductByID() {
	s.Require().NotNil(s.dbContext)

	id := s.items[0].ID

	exist := gormdbcontext.Exists[*datamodel.ProductDataModel](
		context.Background(),
		s.dbContext,
		id,
	)
	s.Require().True(exist)
}

// TestNoneExistsProductByID tests the none exists product by id.
func (s *DBContextTestSuite) TestNoneExistsProductByID() {
	s.Require().NotNil(s.dbContext)

	id := uuid.NewV4()

	exist := gormdbcontext.Exists[*datamodel.ProductDataModel](
		context.Background(),
		s.dbContext,
		id,
	)

	s.Require().False(exist)
}

// TestDeleteProductByID tests the delete product by id.
func (s *DBContextTestSuite) TestDeleteProductByID() {
	s.Require().NotNil(s.dbContext)

	id := s.items[0].ID

	err := gormdbcontext.DeleteDataModelByID[*datamodel.ProductDataModel](
		context.Background(),
		s.dbContext,
		id,
	)
	s.Require().NoError(err)

	p, err := gormdbcontext.FindModelByID[*datamodel.ProductDataModel, *models.Product](
		context.Background(),
		s.dbContext,
		id,
	)
	s.Require().Error(err)
	s.Require().Nil(p)

	// https://gorm.io/docs/delete.html#Find-soft-deleted-records
	var softDeletedProduct *datamodel.ProductDataModel
	s.dbContext.DB().Scopes(scopes.FilterAllItemsWithSoftDeleted).First(&softDeletedProduct, id)
	s.Require().NotNil(softDeletedProduct)

	var deletedCount int64
	var allCount int64

	// https://gorm.io/docs/advanced_query.html#Count
	s.dbContext.DB().
		Model(&datamodel.ProductDataModel{}).
		Scopes(scopes.FilterAllItemsWithSoftDeleted).
		Count(&allCount)
	s.Equal(allCount, int64(2))

	s.dbContext.DB().
		Model(&datamodel.ProductDataModel{}).
		Scopes(scopes.SoftDeleted).
		Count(&deletedCount)
	s.Equal(deletedCount, int64(1))
}

// TestCreateProduct tests the create product.
func (s *DBContextTestSuite) TestCreateProduct() {
	s.Require().NotNil(s.dbContext)

	item := &models.Product{
		ID:          uuid.NewV4(),
		Name:        gofakeit.Name(),
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(100, 1000),
	}

	res, err := gormdbcontext.AddModel[*datamodel.ProductDataModel, *models.Product](
		context.Background(),
		s.dbContext,
		item,
	)
	s.Require().NoError(err)

	p, err := gormdbcontext.FindModelByID[*datamodel.ProductDataModel, *models.Product](
		context.Background(),
		s.dbContext,
		item.ID,
	)
	s.Require().NoError(err)
	s.Require().NotNil(p)

	s.Assert().Equal(p.ID, item.ID)
	s.Assert().Equal(p.ID, res.ID)
}

// TestUpdateProduct tests the update product.
func (s *DBContextTestSuite) TestUpdateProduct() {
	s.Require().NotNil(s.dbContext)

	id := s.items[0].ID

	p, err := gormdbcontext.FindModelByID[*datamodel.ProductDataModel, *models.Product](
		context.Background(),
		s.dbContext,
		id,
	)
	s.Require().NoError(err)

	newName := gofakeit.Name()
	item := p
	item.Name = newName

	res, err := gormdbcontext.UpdateModel[*datamodel.ProductDataModel, *models.Product](
		context.Background(),
		s.dbContext,
		item,
	)
	s.Require().NoError(err)

	p2, err := gormdbcontext.FindModelByID[*datamodel.ProductDataModel, *models.Product](
		context.Background(),
		s.dbContext,
		id,
	)
	s.Require().NoError(err)

	s.Assert().Equal(item.Name, p2.Name)
	s.Assert().Equal(res.Name, p2.Name)
}

// SetupTest sets up the test.
func (s *DBContextTestSuite) SetupTest() {
	err := mappings.ConfigureProductsMappings()
	s.Require().NoError(err)

	var gormDBContext *CatalogsGormDBContext
	var gormOptions *gormPostgres.GormOptions

	app := fxtest.New(
		s.T(),
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		gormPostgres.Module,
		fx.Decorate(
			func(cfg *gormPostgres.GormOptions) (*gormPostgres.GormOptions, error) {
				// using sql-lite with a database file
				cfg.UseSQLLite = true

				return cfg, nil
			},
		),
		fx.Provide(NewCatalogsDBContext),
		fx.Populate(&gormDBContext),
		fx.Populate(&gormOptions),
	).RequireStart()

	s.app = app
	s.dbContext = gormDBContext
	s.dbFilePath = gormOptions.DND()

	s.initDB()
}

func (s *DBContextTestSuite) TearDownTest() {
	err := s.cleanupDB()
	s.Require().NoError(err)

	mapper.ClearMappings()

	s.app.RequireStop()
}

func (s *DBContextTestSuite) initDB() {
	err := migrateGorm(s.dbContext.DB())
	s.Require().NoError(err)

	products, err := seedData(s.dbContext.DB())
	s.Require().NoError(err)

	s.items = products
}

func (s *DBContextTestSuite) cleanupDB() error {
	sqldb, _ := s.dbContext.DB().DB()
	e := sqldb.Close()
	s.Require().NoError(e)

	// removing sql-lite file
	err := os.Remove(s.dbFilePath)

	return err
}

func migrateGorm(db *gorm.DB) error {
	err := db.AutoMigrate(&datamodel.ProductDataModel{})
	if err != nil {
		return err
	}

	return nil
}

func seedData(gormDB *gorm.DB) ([]*datamodel.ProductDataModel, error) {
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
	err := gormDB.CreateInBatches(products, len(products)).Error
	if err != nil {
		return nil, errors.Wrap(err, "error in seed database")
	}

	return products, nil
}
