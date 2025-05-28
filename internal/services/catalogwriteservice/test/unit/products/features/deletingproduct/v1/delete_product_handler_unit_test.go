//go:build unit
// +build unit

package v1

import (
	"net/http"
	"testing"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/gormdbcontext"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/data/datamodels"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	deletingproductv1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/deletingproduct/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"
)

type deleteProductHandlerUnitTests struct {
	*unittest.CatalogWriteUnitTestSharedFixture
	handler cqrs.RequestHandlerWithRegisterer[*deletingproductv1.DeleteProduct, *mediatr.Unit]
}

func TestDeleteProductHandlerUnit(t *testing.T) {
	suite.Run(
		t,
		&deleteProductHandlerUnitTests{
			CatalogWriteUnitTestSharedFixture: unittest.NewCatalogWriteUnitTestSharedFixture(t),
		},
	)
}

func (c *deleteProductHandlerUnitTests) SetupTest() {
	// call base SetupTest hook before running child hook
	c.CatalogWriteUnitTestSharedFixture.SetupTest()
	c.handler = deletingproductv1.NewDeleteProductHandler(
		fxparams.ProductHandlerParams{
			Log:               c.Log,
			CatalogsDBContext: c.CatalogDBContext,
			RabbitmqProducer:  c.Bus,
			Tracer:            c.Tracer,
		},
	)
}

func (c *deleteProductHandlerUnitTests) TearDownTest() {
	// call base TearDownTest hook before running child hook
	c.CatalogWriteUnitTestSharedFixture.TearDownTest()
}

func (c *deleteProductHandlerUnitTests) Test_Handle_Should_Delete_Product_With_Valid_Id() {
	id := c.Products[0].ID

	deleteProduct := &deletingproductv1.DeleteProduct{
		ProductID: id,
	}

	c.BeginTx()
	_, err := c.handler.Handle(c.Ctx, deleteProduct)
	c.CommitTx()

	c.Require().NoError(err)

	p, err := gormdbcontext.FindDataModelByID[*datamodels.ProductDataModel](
		c.Ctx,
		c.CatalogDBContext,
		id,
	)

	c.Require().Nil(p)
	c.Require().Error(err)

	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
}

func (c *deleteProductHandlerUnitTests) Test_Handle_Should_Return_NotFound_Error_When_Id_Is_Invalid() {
	id := uuid.NewV4()

	deleteProduct := &deletingproductv1.DeleteProduct{
		ProductID: id,
	}

	c.BeginTx()
	res, err := c.handler.Handle(c.Ctx, deleteProduct)
	c.CommitTx()

	c.Require().Error(err)
	c.True(customErrors.IsNotFoundError(err))
	c.Nil(res)

	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 0)
}

func (c *deleteProductHandlerUnitTests) Test_Handle_Should_Return_Error_For_Error_In_Bus() {
	id := c.Products[0].ID

	deleteProduct := &deletingproductv1.DeleteProduct{
		ProductID: id,
	}

	// override called mock
	// https://github.com/stretchr/testify/issues/558
	c.Bus.Mock.ExpectedCalls = nil
	c.Bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return(errors.New("error in the publish message"))

	c.BeginTx()
	dto, err := c.handler.Handle(c.Ctx, deleteProduct)
	c.CommitTx()

	c.Nil(dto)

	c.Bus.AssertNumberOfCalls(c.T(), "PublishMessage", 1)
	c.True(customErrors.IsApplicationError(err, http.StatusInternalServerError))
	c.ErrorContains(err, "error in publishing 'ProductDeleted' message")
}
