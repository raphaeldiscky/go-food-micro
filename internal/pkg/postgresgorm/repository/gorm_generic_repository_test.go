//go:build integration
// +build integration

// Package repository provides the gorm generic repository.
package repository

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	_ "github.com/lib/pq" // postgres driver

	gofakeit "github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/data"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	defaultLogger "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm"
	gorm2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/testcontainer/gorm"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
)

const (
	TestUpdatedProductName = "product2_updated"
)

// Product is a domain_events entity.
type Product struct {
	ID          uuid.UUID
	Name        string
	Weight      int
	IsAvailable bool
}

// ProductGorm is DTO used to map Product entity to database.
type ProductGorm struct {
	ID          uuid.UUID `gorm:"primaryKey;column:id"`
	Name        string    `gorm:"column:name"`
	Weight      int       `gorm:"column:weight"`
	IsAvailable bool      `gorm:"column:is_available"`
}

// TableName returns the table name.
func (v *ProductGorm) TableName() string {
	return "products_gorm"
}

// GormGenericRepositoryTestSuite is the test suite for the gorm generic repository.
type gormGenericRepositoryTest struct {
	suite.Suite
	DB                             *gorm.DB
	productRepository              data.GenericRepository[*ProductGorm]
	productRepositoryWithDataModel data.GenericRepositoryWithDataModel[*ProductGorm, *Product]
	products                       []*ProductGorm
}

// TestGormGenericRepository tests the gorm generic repository.
func TestGormGenericRepository(t *testing.T) {
	suite.Run(
		t,
		&gormGenericRepositoryTest{},
	)
}

// SetupSuite sets up the test suite.
func (c *gormGenericRepositoryTest) SetupSuite() {
	opts, err := gorm2.NewGormTestContainers(defaultLogger.GetLogger()).
		PopulateContainerOptions(context.Background(), c.T())
	c.Require().NoError(err)

	gormDB, err := postgresgorm.NewGorm(opts)
	c.Require().NoError(err)
	c.DB = gormDB

	err = migrationDatabase(gormDB)
	c.Require().NoError(err)

	c.productRepository = NewGenericGormRepository[*ProductGorm](gormDB)
	c.productRepositoryWithDataModel = NewGenericGormRepositoryWithDataModel[*ProductGorm, *Product](
		gormDB,
	)

	err = mapper.CreateMap[*ProductGorm, *Product]()
	if err != nil {
		log.Fatal(err)
	}

	err = mapper.CreateMap[*Product, *ProductGorm]()
	if err != nil {
		log.Fatal(err)
	}
}

// SetupTest sets up the test.
func (c *gormGenericRepositoryTest) SetupTest() {
	p, err := seedData(context.Background(), c.DB)
	c.Require().NoError(err)
	c.products = p
}

// TearDownTest tears down the test.
func (c *gormGenericRepositoryTest) TearDownTest() {
	err := c.cleanupPostgresData()
	c.Require().NoError(err)
}

// TestAdd tests the add.
func (c *gormGenericRepositoryTest) TestAdd() {
	ctx := context.Background()

	product := &ProductGorm{
		ID:          uuid.NewV4(),
		Name:        gofakeit.Name(),
		Weight:      gofakeit.Number(100, 1000),
		IsAvailable: true,
	}

	err := c.productRepository.Add(ctx, product)
	c.Require().NoError(err)

	p, err := c.productRepository.GetByID(ctx, product.ID)
	if err != nil {
		return
	}

	c.Assert().NotNil(p)
	c.Assert().Equal(product.ID, p.ID)
}

// TestAddWithDataModel tests the add with data model.
func (c *gormGenericRepositoryTest) TestAddWithDataModel() {
	ctx := context.Background()

	product := &Product{
		ID:          uuid.NewV4(),
		Name:        gofakeit.Name(),
		Weight:      gofakeit.Number(100, 1000),
		IsAvailable: true,
	}

	err := c.productRepositoryWithDataModel.Add(ctx, product)
	c.Require().NoError(err)

	p, err := c.productRepositoryWithDataModel.GetByID(ctx, product.ID)
	if err != nil {
		return
	}

	c.Assert().NotNil(p)
	c.Assert().Equal(product.ID, p.ID)
}

// TestGetById tests the get by id.
func (c *gormGenericRepositoryTest) TestGetById() {
	ctx := context.Background()

	all, err := c.productRepository.GetAll(ctx, utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	p := all.Items[0]

	testCases := []struct {
		Name         string
		ProductID    uuid.UUID
		ExpectResult *ProductGorm
	}{
		{
			Name:         p.Name,
			ProductID:    p.ID,
			ExpectResult: p,
		},
		{
			Name:         "NonExistingProduct",
			ProductID:    uuid.NewV4(),
			ExpectResult: nil,
		},
	}

	for _, s := range testCases {
		c.T().Run(s.Name, func(t *testing.T) {
			// Remove t.Parallel() to avoid race conditions with shared database state
			res, err := c.productRepository.GetByID(ctx, s.ProductID)
			if s.ExpectResult == nil {
				assert.Error(t, err)
				assert.True(t, customErrors.IsNotFoundError(err))
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				if res != nil { // Fix nil pointer dereference
					assert.Equal(t, s.ExpectResult.ID, res.ID)
				}
			}
		})
	}
}

// TestGetByIdWithDataModel tests the get by id with data model.
func (c *gormGenericRepositoryTest) TestGetByIdWithDataModel() {
	ctx := context.Background()

	all, err := c.productRepositoryWithDataModel.GetAll(
		ctx,
		utils.NewListQuery(10, 1),
	)
	if err != nil {
		return
	}
	p := all.Items[0]

	testCases := []struct {
		Name         string
		ProductID    uuid.UUID
		ExpectResult *Product
	}{
		{
			Name:         p.Name,
			ProductID:    p.ID,
			ExpectResult: p,
		},
		{
			Name:         "NonExistingProduct",
			ProductID:    uuid.NewV4(),
			ExpectResult: nil,
		},
	}

	for _, s := range testCases {
		c.T().Run(s.Name, func(t *testing.T) {
			// Remove t.Parallel() to avoid race conditions with shared database state
			res, err := c.productRepositoryWithDataModel.GetByID(
				ctx,
				s.ProductID,
			)

			if s.ExpectResult == nil {
				assert.Error(t, err)
				assert.True(t, customErrors.IsNotFoundError(err))
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				if res != nil { // Fix nil pointer dereference
					assert.Equal(t, s.ExpectResult.ID, res.ID)
				}
			}
		})
	}
}

// TestGetAll tests the get all.
func (c *gormGenericRepositoryTest) TestGetAll() {
	ctx := context.Background()

	models, err := c.productRepository.GetAll(ctx, utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	c.Assert().NotEmpty(models.Items)
}

// TestGetAllWithDataModel tests the get all with data model.
func (c *gormGenericRepositoryTest) TestGetAllWithDataModel() {
	ctx := context.Background()

	models, err := c.productRepositoryWithDataModel.GetAll(
		ctx,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	c.Assert().NotEmpty(models.Items)
}

// TestSearch tests the search.
func (c *gormGenericRepositoryTest) TestSearch() {
	ctx := context.Background()

	models, err := c.productRepository.Search(
		ctx,
		c.products[0].Name,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	c.Assert().NotEmpty(models.Items)
	c.Assert().Equal(len(models.Items), 1)
}

// TestSearchWithDataModel tests the search with data model.
func (c *gormGenericRepositoryTest) TestSearchWithDataModel() {
	ctx := context.Background()

	models, err := c.productRepositoryWithDataModel.Search(
		ctx,
		c.products[0].Name,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	c.Assert().NotEmpty(models.Items)
	c.Assert().Equal(len(models.Items), 1)
}

// TestWhere tests the where.
func (c *gormGenericRepositoryTest) TestWhere() {
	ctx := context.Background()

	models, err := c.productRepository.GetByFilter(
		ctx,
		map[string]interface{}{"name": c.products[0].Name},
	)
	c.Require().NoError(err)

	c.Assert().NotEmpty(models)
	c.Assert().Equal(len(models), 1)
}

// TestWhereWithDataModel tests the where with data model.
func (c *gormGenericRepositoryTest) TestWhereWithDataModel() {
	ctx := context.Background()

	models, err := c.productRepositoryWithDataModel.GetByFilter(
		ctx,
		map[string]interface{}{"name": c.products[0].Name},
	)
	c.Require().NoError(err)

	c.Assert().NotEmpty(models)
	c.Assert().Equal(len(models), 1)
}

// TestUpdate tests the update.
func (c *gormGenericRepositoryTest) TestUpdate() {
	ctx := context.Background()

	products, err := c.productRepository.GetAll(ctx, utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	product := products.Items[0]

	product.Name = TestUpdatedProductName
	err = c.productRepository.Update(ctx, product)
	c.Require().NoError(err)

	single, err := c.productRepository.GetByID(ctx, product.ID)
	c.Require().NoError(err)

	c.Assert().NotNil(single)
	c.Assert().Equal(TestUpdatedProductName, single.Name)
}

// TestUpdateWithDataModel tests the update with data model.
func (c *gormGenericRepositoryTest) TestUpdateWithDataModel() {
	ctx := context.Background()

	products, err := c.productRepositoryWithDataModel.GetAll(
		ctx,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	product := products.Items[0]

	product.Name = TestUpdatedProductName
	err = c.productRepositoryWithDataModel.Update(ctx, product)
	c.Require().NoError(err)

	single, err := c.productRepositoryWithDataModel.GetByID(ctx, product.ID)
	c.Require().NoError(err)

	c.Assert().NotNil(single)
	c.Assert().Equal(TestUpdatedProductName, single.Name)
}

// cleanupPostgresData cleans up the postgres data.
func (c *gormGenericRepositoryTest) cleanupPostgresData() error {
	tables := []string{"products_gorm"}
	// Iterate over the tables and delete all records
	for _, table := range tables {
		err := c.DB.Exec("DELETE FROM " + table).Error

		return err
	}

	return nil
}

// migrationDatabase migrates the database.
func migrationDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(ProductGorm{})
	if err != nil {
		return err
	}

	return nil
}

// seedData seeds the data.
func seedData(ctx context.Context, db *gorm.DB) ([]*ProductGorm, error) {
	seedProducts := []*ProductGorm{
		{
			ID:          uuid.NewV4(),
			Name:        "seed_product1",
			Weight:      100,
			IsAvailable: true,
		},
		{
			ID:          uuid.NewV4(),
			Name:        "seed_product2",
			Weight:      100,
			IsAvailable: true,
		},
	}

	err := db.WithContext(ctx).Create(seedProducts).Error
	if err != nil {
		return nil, err
	}

	return seedProducts, nil
}
