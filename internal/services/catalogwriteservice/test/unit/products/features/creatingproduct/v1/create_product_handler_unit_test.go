//go:build unit
// +build unit

package v1

import (
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/gormdbcontext"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	uuid "github.com/satori/go.uuid"

	datamodels "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/data/datamodels"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	creatingproductv1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1"
	creatingproductdtosv1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/models"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"
)

type createProductHandlerUnitTests struct {
	*unittest.CatalogWriteUnitTestSharedFixture
	handler cqrs.RequestHandlerWithRegisterer[*creatingproductv1.CreateProduct, *creatingproductdtosv1.CreateProductResponseDto]
}

func TestCreateProductHandlerUnit(t *testing.T) {
	suite.Run(t, &createProductHandlerUnitTests{
		CatalogWriteUnitTestSharedFixture: unittest.NewCatalogWriteUnitTestSharedFixture(t),
	},
	)
}

func (c *createProductHandlerUnitTests) SetupTest() {
	// call base SetupTest hook before running child hook
	c.CatalogWriteUnitTestSharedFixture.SetupTest()
	c.handler = creatingproductv1.NewCreateProductHandler(
		fxparams.ProductHandlerParams{
			CatalogsDBContext: c.CatalogDBContext,
			Tracer:            c.Tracer,
			RabbitmqProducer:  c.Bus,
			Log:               c.Log,
		},
	)
}

func (c *createProductHandlerUnitTests) TearDownTest() {
	// call base TearDownTest hook before running child hook
	c.CatalogWriteUnitTestSharedFixture.TearDownTest()
}

// TestHandleShouldCreateNewProductWithValidData tests the handle should create new product with valid data.
func (c *createProductHandlerUnitTests) TestHandleShouldCreateNewProductWithValidData() {
	id := uuid.NewV4()

	createProduct := &creatingproductv1.CreateProduct{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	c.BeginTx()
	_, err := c.handler.Handle(c.Ctx, createProduct)
	c.CommitTx()

	c.Require().NoError(err)

	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)

	res, err := gormdbcontext.FindModelByID[*datamodels.ProductDataModel, *models.Product](
		c.Ctx,
		c.CatalogDBContext,
		id,
	)
	c.Require().NoError(err)

	c.Assert().Equal(res.ID, createProduct.ProductID)
}

// TestHandleShouldReturnErrorForDuplicateItem tests the handle should return error for duplicate item.
func (c *createProductHandlerUnitTests) TestHandleShouldReturnErrorForDuplicateItem() {
	id := uuid.NewV4()

	createProduct := &creatingproductv1.CreateProduct{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	c.BeginTx()
	dto, err := c.handler.Handle(c.Ctx, createProduct)
	c.Require().NoError(err)
	c.Require().NotNil(dto)
	c.CommitTx()

	c.BeginTx()
	dto, err = c.handler.Handle(c.Ctx, createProduct)
	c.CommitTx()

	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.True(customErrors.IsConflictError(err))
	c.ErrorContains(err, "product already exists")
	c.Nil(dto)
}

// TestHandleShouldReturnErrorForErrorInBus tests the handle should return error for error in bus.
func (c *createProductHandlerUnitTests) TestHandleShouldReturnErrorForErrorInBus() {
	id := uuid.NewV4()

	createProduct := &creatingproductv1.CreateProduct{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	// override called mock
	// https://github.com/stretchr/testify/issues/558
	c.Bus.Mock.ExpectedCalls = nil
	c.Bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return(errors.New("error in the publish message"))

	c.BeginTx()

	dto, err := c.handler.Handle(c.Ctx, createProduct)

	c.CommitTx()

	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.ErrorContains(err, "error in the publish message")
	c.ErrorContains(
		err,
		"error in publishing ProductCreated integration_events event",
	)
	c.Nil(dto)
}

// TestHandleShouldReturnErrorForErrorInMapping tests the handle should return error for error in mapping.
func (c *createProductHandlerUnitTests) TestHandleShouldReturnErrorForErrorInMapping() {
	id := uuid.NewV4()

	createProduct := &creatingproductv1.CreateProduct{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	mapper.ClearMappings()

	c.BeginTx()

	dto, err := c.handler.Handle(c.Ctx, createProduct)

	c.CommitTx()

	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 0)
	c.ErrorContains(err, "error in the mapping")
	c.True(customErrors.IsInternalServerError(err))
	c.Nil(dto)
}
