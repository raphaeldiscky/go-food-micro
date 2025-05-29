//go:build unit
// +build unit

package v1

import (
	"fmt"
	"net/http"
	"testing"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/gormdbcontext"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/data/datamodels"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	updatingoroductsv1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/updatingproduct/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"
)

type updateProductHandlerUnitTests struct {
	*unittest.CatalogWriteUnitTestSharedFixture
	handler cqrs.RequestHandlerWithRegisterer[*updatingoroductsv1.UpdateProduct, *mediatr.Unit]
}

func TestUpdateProductHandlerUnit(t *testing.T) {
	suite.Run(
		t,
		&updateProductHandlerUnitTests{
			CatalogWriteUnitTestSharedFixture: unittest.NewCatalogWriteUnitTestSharedFixture(t),
		},
	)
}

func (c *updateProductHandlerUnitTests) SetupTest() {
	// call base `SetupTest hook` before running child hook
	c.CatalogWriteUnitTestSharedFixture.SetupTest()
	c.handler = updatingoroductsv1.NewUpdateProductHandler(
		fxparams.ProductHandlerParams{
			CatalogsDBContext: c.CatalogDBContext,
			Tracer:            c.Tracer,
			RabbitmqProducer:  c.Bus,
			Log:               c.Log,
		},
	)
}

func (c *updateProductHandlerUnitTests) TearDownTest() {
	// call base `TearDownTest hook` before running child hook
	c.CatalogWriteUnitTestSharedFixture.TearDownTest()
}

// TestHandleShouldUpdateProductWithValidData tests the handle should update product with valid data.
func (c *updateProductHandlerUnitTests) TestHandleShouldUpdateProductWithValidData() {
	existing := c.Products[0]

	updateProductCommand, err := updatingoroductsv1.NewUpdateProductWithValidation(
		existing.ID,
		gofakeit.Name(),
		gofakeit.EmojiDescription(),
		existing.Price,
	)
	c.Require().NoError(err)

	c.BeginTx()
	_, err = c.handler.Handle(c.Ctx, updateProductCommand)
	c.CommitTx()

	c.Require().NoError(err)

	updatedProduct, err := gormdbcontext.FindDataModelByID[*datamodels.ProductDataModel](
		c.Ctx,
		c.CatalogDBContext,
		updateProductCommand.ProductID,
	)
	c.Require().NoError(err)

	c.Assert().Equal(updatedProduct.ID, updateProductCommand.ProductID)
	c.Assert().Equal(updatedProduct.Name, updateProductCommand.Name)
	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
}

// TestHandleShouldReturnErrorForNotFoundItem tests the handle should return error for not found item.
func (c *updateProductHandlerUnitTests) TestHandleShouldReturnErrorForNotFoundItem() {
	id := uuid.NewV4()

	command, err := updatingoroductsv1.NewUpdateProductWithValidation(
		id,
		gofakeit.Name(),
		gofakeit.EmojiDescription(),
		gofakeit.Price(150, 6000),
	)
	c.Require().NoError(err)

	c.BeginTx()
	_, err = c.handler.Handle(c.Ctx, command)
	c.CommitTx()

	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 0)
	c.True(customErrors.IsApplicationError(err, http.StatusNotFound))
	c.ErrorContains(
		err,
		fmt.Sprintf("product with id `%s` not found", id.String()),
	)
}

// TestHandleShouldReturnErrorForErrorInBus tests the handle should return error for error in bus.
func (c *updateProductHandlerUnitTests) TestHandleShouldReturnErrorForErrorInBus() {
	existing := c.Products[0]

	updateProductCommand, err := updatingoroductsv1.NewUpdateProductWithValidation(
		existing.ID,
		gofakeit.Name(),
		gofakeit.EmojiDescription(),
		existing.Price,
	)
	c.Require().NoError(err)

	// override called mock
	// https://github.com/stretchr/testify/issues/558
	c.Bus.Mock.ExpectedCalls = nil
	c.Bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return(errors.New("error in the publish message"))

	c.BeginTx()
	_, err = c.handler.Handle(c.Ctx, updateProductCommand)
	c.CommitTx()

	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.ErrorContains(err, "error in the publish message")
	c.ErrorContains(err, "error in publishing 'ProductUpdated' message")
}
