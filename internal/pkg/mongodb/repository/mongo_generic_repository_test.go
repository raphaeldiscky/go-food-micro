//go:build integration
// +build integration

// Package repository provides the mongodb generic repository.
package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/data"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	defaultLogger "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
	mongocontainer "github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/testcontainer/mongo"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
)

const (
	TestUpdatedProductName = "product2_updated"
)

// Product is a domain_events entity.
type Product struct {
	ID          string
	Name        string
	Weight      int
	IsAvailable bool
}

// ProductMongo is the model for the product.
type ProductMongo struct {
	ID          string `json:"id"          bson:"_id,omitempty"` // https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/write-operations/insert/#the-_id-field
	Name        string `json:"name"        bson:"name"`
	Weight      int    `json:"weight"      bson:"weight"`
	IsAvailable bool   `json:"isAvailable" bson:"isAvailable"`
}

// mongoGenericRepositoryTest is the test suite for the mongodb generic repository.
type mongoGenericRepositoryTest struct {
	suite.Suite
	databaseName                   string
	collectionName                 string
	mongoClient                    *mongo.Client
	productRepository              data.GenericRepository[*ProductMongo]
	productRepositoryWithDataModel data.GenericRepositoryWithDataModel[*ProductMongo, *Product]
	products                       []*ProductMongo
}

// TestMongoGenericRepository tests the mongodb generic repository.
func TestMongoGenericRepository(t *testing.T) {
	t.Helper()
	suite.Run(
		t,
		&mongoGenericRepositoryTest{
			databaseName:   "catalogs_write",
			collectionName: "products",
		},
	)
}

// SetupSuite sets up the test suite.
func (c *mongoGenericRepositoryTest) SetupSuite() {
	opts, err := mongocontainer.NewMongoTestContainers(defaultLogger.GetLogger()).
		PopulateContainerOptions(context.Background(), c.T())
	c.Require().NoError(err)

	mongoClient, err := mongodb.NewMongoDB(opts)
	c.Require().NoError(err)
	c.mongoClient = mongoClient

	c.productRepository = NewGenericMongoRepository[*ProductMongo](
		mongoClient,
		c.databaseName,
		c.collectionName,
	)
	c.productRepositoryWithDataModel = NewGenericMongoRepositoryWithDataModel[*ProductMongo, *Product](
		mongoClient,
		c.databaseName,
		c.collectionName,
	)

	err = mapper.CreateMap[*ProductMongo, *Product]()
	c.Require().NoError(err)

	err = mapper.CreateMap[*Product, *ProductMongo]()
	c.Require().NoError(err)
}

// SetupTest sets up the test.
func (c *mongoGenericRepositoryTest) SetupTest() {
	p, err := c.seedData(context.Background())
	c.Require().NoError(err)
	c.products = p
}

// TearDownTest tears down the test.
func (c *mongoGenericRepositoryTest) TearDownTest() {
	err := c.cleanupMongoData()
	c.Require().NoError(err)
}

// TestAdd tests the add.
func (c *mongoGenericRepositoryTest) TestAdd() {
	ctx := context.Background()

	product := &ProductMongo{
		// we generate id ourselves because auto generate mongo string id column with type _id is not an uuid
		ID:          uuid.NewV4().String(),
		Name:        gofakeit.Name(),
		Weight:      gofakeit.Number(100, 1000),
		IsAvailable: true,
	}

	err := c.productRepository.Add(ctx, product)
	c.Require().NoError(err)

	id, err := uuid.FromString(product.ID)
	c.Require().NoError(err)

	p, err := c.productRepository.GetByID(ctx, id)
	c.Require().NoError(err)

	c.Assert().NotNil(p)
	c.Assert().Equal(product.ID, p.ID)
}

// TestAddWithDataModel tests the add with data model.
func (c *mongoGenericRepositoryTest) TestAddWithDataModel() {
	ctx := context.Background()

	product := &ProductMongo{
		// we generate id ourselves because auto generate mongo string id column with type _id is not an uuid
		ID:          uuid.NewV4().String(),
		Name:        gofakeit.Name(),
		Weight:      gofakeit.Number(100, 1000),
		IsAvailable: true,
	}

	err := c.productRepository.Add(ctx, product)
	c.Require().NoError(err)

	id, err := uuid.FromString(product.ID)
	c.Require().NoError(err)

	p, err := c.productRepository.GetByID(ctx, id)
	c.Require().NoError(err)

	c.Assert().NotNil(p)
	c.Assert().Equal(product.ID, p.ID)
}

// TestGetById tests the get by id.
func (c *mongoGenericRepositoryTest) TestGetById() {
	ctx := context.Background()

	all, err := c.productRepository.GetAll(ctx, utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	p := all.Items[0]
	id, err := uuid.FromString(p.ID)
	c.Require().NoError(err)
	name := p.Name

	testCases := []struct {
		Name         string
		ProductID    uuid.UUID
		ExpectResult *ProductMongo
	}{
		{
			Name:         name,
			ProductID:    id,
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
			// Remove t.Parallel() to avoid race conditions with shared database
			res, err := c.productRepository.GetByID(ctx, s.ProductID)
			if s.ExpectResult == nil {
				assert.Error(t, err)
				assert.True(t, customErrors.IsNotFoundError(err))
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, s.ExpectResult.ID, res.ID)
			}
		})
	}
}

// TestGetByIdWithDataModel tests the get by id with data model.
func (c *mongoGenericRepositoryTest) TestGetByIdWithDataModel() {
	ctx := context.Background()

	all, err := c.productRepositoryWithDataModel.GetAll(
		ctx,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	p := all.Items[0]
	id, err := uuid.FromString(p.ID)
	c.Require().NoError(err)
	name := p.Name

	testCases := []struct {
		Name         string
		ProductID    uuid.UUID
		ExpectResult *Product
	}{
		{
			Name:         name,
			ProductID:    id,
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
			// Remove t.Parallel() to avoid race conditions with shared database
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
				assert.Equal(t, s.ExpectResult.ID, res.ID)
			}
		})
	}
}

// TestFirstOrDefault tests the first or default.
func (c *mongoGenericRepositoryTest) TestFirstOrDefault() {
	ctx := context.Background()

	all, err := c.productRepository.GetAll(ctx, utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	p := all.Items[0]

	single, err := c.productRepository.FirstOrDefault(
		ctx,
		map[string]interface{}{"_id": p.ID},
	)
	c.Require().NoError(err)
	c.Assert().NotNil(single)
}

// TestFirstOrDefaultWithDataModel tests the first or default with data model.
func (c *mongoGenericRepositoryTest) TestFirstOrDefaultWithDataModel() {
	ctx := context.Background()

	all, err := c.productRepositoryWithDataModel.GetAll(
		ctx,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	p := all.Items[0]

	single, err := c.productRepositoryWithDataModel.FirstOrDefault(
		ctx,
		map[string]interface{}{"_id": p.ID},
	)

	c.Require().NoError(err)
	c.Assert().NotNil(single)
}

// TestGetAll tests the get all.
func (c *mongoGenericRepositoryTest) TestGetAll() {
	ctx := context.Background()

	models, err := c.productRepository.GetAll(ctx, utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	c.Assert().NotEmpty(models.Items)
}

// TestGetAllWithDataModel tests the get all with data model.
func (c *mongoGenericRepositoryTest) TestGetAllWithDataModel() {
	ctx := context.Background()

	models, err := c.productRepositoryWithDataModel.GetAll(
		ctx,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	c.Assert().NotEmpty(models.Items)
}

// TestSearch tests the search.
func (c *mongoGenericRepositoryTest) TestSearch() {
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
func (c *mongoGenericRepositoryTest) TestSearchWithDataModel() {
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

// TestGetByFilter tests the get by filter.
func (c *mongoGenericRepositoryTest) TestGetByFilter() {
	ctx := context.Background()

	models, err := c.productRepository.GetByFilter(
		ctx,
		map[string]interface{}{"name": c.products[0].Name},
	)
	c.Require().NoError(err)

	c.Assert().NotEmpty(models)
	c.Assert().Equal(len(models), 1)
}

// TestGetByFilterWithDataModel tests the get by filter with data model.
func (c *mongoGenericRepositoryTest) TestGetByFilterWithDataModel() {
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
func (c *mongoGenericRepositoryTest) TestUpdate() {
	ctx := context.Background()

	products, err := c.productRepository.GetAll(ctx, utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	product := products.Items[0]

	product.Name = TestUpdatedProductName
	err = c.productRepository.Update(ctx, product)
	c.Require().NoError(err)

	id, err := uuid.FromString(product.ID)
	c.Require().NoError(err)

	single, err := c.productRepository.GetByID(ctx, id)
	c.Require().NoError(err)

	c.Assert().NotNil(single)
	c.Assert().Equal(TestUpdatedProductName, single.Name)
}

// TestUpdateWithDataModel tests the update with data model.
func (c *mongoGenericRepositoryTest) TestUpdateWithDataModel() {
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

	id, err := uuid.FromString(product.ID)
	c.Require().NoError(err)

	single, err := c.productRepositoryWithDataModel.GetByID(ctx, id)
	c.Require().NoError(err)

	c.Assert().NotNil(single)
	c.Assert().Equal(TestUpdatedProductName, single.Name)
}

// TestDelete tests the delete.
func (c *mongoGenericRepositoryTest) TestDelete() {
	ctx := context.Background()

	products, err := c.productRepository.GetAll(ctx, utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	product := products.Items[0]

	id, err := uuid.FromString(product.ID)
	c.Require().NoError(err)

	err = c.productRepository.Delete(ctx, id)
	c.Require().NoError(err)

	single, err := c.productRepository.GetByID(ctx, id)
	c.Require().Error(err)
	c.Assert().True(customErrors.IsNotFoundError(err))
	c.Assert().Nil(single)
}

// cleanupMongoData cleans up the mongo data.
func (c *mongoGenericRepositoryTest) cleanupMongoData() error {
	collections := []string{c.collectionName}
	err := cleanupCollections(
		c.mongoClient,
		collections,
		c.databaseName,
	)

	return err
}

// cleanupCollections cleans up the collections.
func cleanupCollections(
	db *mongo.Client,
	collections []string,
	databaseName string,
) error {
	database := db.Database(databaseName)
	ctx := context.Background()

	// Iterate over the collections and delete all collections
	for _, collection := range collections {
		collection := database.Collection(collection)

		err := collection.Drop(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// seedData seeds the data.
func (c *mongoGenericRepositoryTest) seedData(
	ctx context.Context,
) ([]*ProductMongo, error) {
	seedProducts := []*ProductMongo{
		{
			ID: uuid.NewV4().
				String(),
			// we generate id ourselves because auto generate mongo string id column with type _id is not an uuid
			Name:        "seed_product1",
			Weight:      100,
			IsAvailable: true,
		},
		{
			ID: uuid.NewV4().
				String(),
			// we generate id ourselves because auto generate mongo string id column with type _id is not an uuid
			Name:        "seed_product2",
			Weight:      100,
			IsAvailable: true,
		},
	}

	// https://go.dev/doc/faq#convert_slice_of_interface
	data := make([]interface{}, len(seedProducts))
	for i, v := range seedProducts {
		data[i] = v
	}

	collection := c.mongoClient.Database(c.databaseName).
		Collection(c.collectionName)
	_, err := collection.InsertMany(ctx, data, &options.InsertManyOptions{})
	if err != nil {
		return nil, err
	}

	return seedProducts, nil
}
