//go:build unit
// +build unit

package v1

import (
	"net/http"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
	"github.com/stretchr/testify/suite"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	gettingproductsv1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/gettingproducts/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/gettingproducts/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"
)

type getProductsHandlerUnitTests struct {
	*unittest.CatalogWriteUnitTestSharedFixture
	handler cqrs.RequestHandlerWithRegisterer[*gettingproductsv1.GetProducts, *dtos.GetProductsResponseDto]
}

func TestGetProductsUnit(t *testing.T) {
	suite.Run(
		t,
		&getProductsHandlerUnitTests{
			CatalogWriteUnitTestSharedFixture: unittest.NewCatalogWriteUnitTestSharedFixture(t),
		},
	)
}

func (c *getProductsHandlerUnitTests) SetupTest() {
	// call base SetupTest hook before running child hook
	c.CatalogWriteUnitTestSharedFixture.SetupTest()
	c.handler = gettingproductsv1.NewGetProductsHandler(
		fxparams.ProductHandlerParams{
			CatalogsDBContext: c.CatalogDBContext,
			Tracer:            c.Tracer,
			RabbitmqProducer:  c.Bus,
			Log:               c.Log,
		})
}

func (c *getProductsHandlerUnitTests) TearDownTest() {
	// call base TearDownTest hook before running child hook
	c.CatalogWriteUnitTestSharedFixture.TearDownTest()
}

// TestHandleShouldReturnProductsSuccessfully tests the handle should return products successfully.
func (c *getProductsHandlerUnitTests) TestHandleShouldReturnProductsSuccessfully() {
	query, err := gettingproductsv1.NewGetProducts(utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	res, err := c.handler.Handle(c.Ctx, query)
	c.Require().NoError(err)
	c.NotNil(res)
	c.NotEmpty(res.Products)
	c.Equal(len(c.Products), len(res.Products.Items))
}

// TestHandleShouldReturnErrorForMappingListResult tests the handle should return error for mapping list result.
func (c *getProductsHandlerUnitTests) TestHandleShouldReturnErrorForMappingListResult() {
	query, err := gettingproductsv1.NewGetProducts(utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	mapper.ClearMappings()

	res, err := c.handler.Handle(c.Ctx, query)
	c.Require().Error(err)
	c.True(customErrors.IsApplicationError(err, http.StatusInternalServerError))
	c.Nil(res)
}
