//go:build unit
// +build unit

package v1

import (
	"fmt"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/stretchr/testify/suite"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	gettingproductbyidv1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/gettingproductbyid/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/gettingproductbyid/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"
)

type getProductByIdHandlerTest struct {
	*unittest.CatalogWriteUnitTestSharedFixture
	handler cqrs.RequestHandlerWithRegisterer[*gettingproductbyidv1.GetProductByID, *dtos.GetProductByIDResponseDto]
}

func TestGetProductByIdHandlerUnit(t *testing.T) {
	suite.Run(t, &getProductByIdHandlerTest{
		CatalogWriteUnitTestSharedFixture: unittest.NewCatalogWriteUnitTestSharedFixture(t),
	})
}

func (c *getProductByIdHandlerTest) SetupTest() {
	// call base SetupTest hook before running child hook
	c.CatalogWriteUnitTestSharedFixture.SetupTest()
	c.handler = gettingproductbyidv1.NewGetProductByIDHandler(
		fxparams.ProductHandlerParams{
			CatalogsDBContext: c.CatalogDBContext,
			Tracer:            c.Tracer,
			RabbitmqProducer:  c.Bus,
			Log:               c.Log,
		})
}

func (c *getProductByIdHandlerTest) TearDownTest() {
	// call base TearDownTest hook before running child hook
	c.CatalogWriteUnitTestSharedFixture.TearDownTest()
}

// TestHandleShouldReturnCorrectProductById tests the handle should return correct product by id.
func (c *getProductByIdHandlerTest) TestHandleShouldReturnCorrectProductById() {
	product := c.Products[0]

	query, err := gettingproductbyidv1.NewGetProductByIDWithValidation(product.ID)
	c.Require().NoError(err)

	dto, err := c.handler.Handle(c.Ctx, query)
	c.Require().NoError(err)
	c.Assert().NotNil(dto)
	c.Assert().NotNil(dto.Product)
	c.Assert().Equal(dto.Product.ID, product.ID)
	c.Assert().Equal(dto.Product.Name, product.Name)
}

// TestHandleShouldReturnNotFoundErrorForNotFoundItem tests the handle should return not found error for not found item.
func (c *getProductByIdHandlerTest) TestHandleShouldReturnNotFoundErrorForNotFoundItem() {
	id := uuid.NewV4()

	query, err := gettingproductbyidv1.NewGetProductByIDWithValidation(id)
	c.Require().NoError(err)

	dto, err := c.handler.Handle(c.Ctx, query)
	c.Require().Error(err)
	c.True(customErrors.IsNotFoundError(err))
	c.ErrorContains(
		err,
		fmt.Sprintf(
			"product with id `%s` not found in the database",
			id.String(),
		),
	)
	c.Nil(dto)
}

// TestHandleShouldReturnErrorForErrorInMapping tests the handle should return error for error in mapping.
func (c *getProductByIdHandlerTest) TestHandleShouldReturnErrorForErrorInMapping() {
	mapper.ClearMappings()

	product := c.Products[0]

	query, err := gettingproductbyidv1.NewGetProductByIDWithValidation(product.ID)
	c.Require().NoError(err)

	dto, err := c.handler.Handle(c.Ctx, query)

	c.Nil(dto)
	c.Require().Error(err)
	c.True(customErrors.IsInternalServerError(err))
	c.ErrorContains(err, "error in the mapping product")
}
