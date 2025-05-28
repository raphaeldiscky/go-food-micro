package v1

import (
	"net/http"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	searchingproductsv1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/searchingproduct/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/searchingproduct/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"

	"github.com/stretchr/testify/suite"
)

// searchProductsHandlerUnitTests is a struct that contains the search products handler unit tests.
type searchProductsHandlerUnitTests struct {
	*unittest.UnitTestSharedFixture
	handler cqrs.RequestHandlerWithRegisterer[*searchingproductsv1.SearchProducts, *dtos.SearchProductsResponseDto]
}

// TestSearchProductsUnit is a function that tests the search products unit.
func TestSearchProductsUnit(t *testing.T) {
	suite.Run(
		t,
		&searchProductsHandlerUnitTests{
			UnitTestSharedFixture: unittest.NewUnitTestSharedFixture(t),
		},
	)
}

// SetupTest is a method that sets up the test.
func (c *searchProductsHandlerUnitTests) SetupTest() {
	// call base SetupTest hook before running child hook
	c.UnitTestSharedFixture.SetupTest()
	c.handler = searchingproductsv1.NewSearchProductsHandler(
		fxparams.ProductHandlerParams{
			CatalogsDBContext: c.CatalogDBContext,
			Tracer:            c.Tracer,
			RabbitmqProducer:  c.Bus,
			Log:               c.Log,
		})
}

// TearDownTest is a method that tears down the test.
func (c *searchProductsHandlerUnitTests) TearDownTest() {
	// call base TearDownTest hook before running child hook
	c.UnitTestSharedFixture.TearDownTest()
}

// Test_Handle_Should_Return_Products_Successfully is a method that tests the handle method should return products successfully.
func (c *searchProductsHandlerUnitTests) Test_Handle_Should_Return_Products_Successfully() {
	query, err := searchingproductsv1.NewSearchProductsWithValidation(
		c.Products[0].Name,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	res, err := c.handler.Handle(c.Ctx, query)
	c.Require().NoError(err)
	c.NotNil(res)
	c.NotEmpty(res.Products)
	c.Equal(len(res.Products.Items), 1)
}

// Test_Handle_Should_Return_Error_For_Mapping_List_Result is a method that tests the handle method should return error for mapping list result.
func (c *searchProductsHandlerUnitTests) Test_Handle_Should_Return_Error_For_Mapping_List_Result() {
	query, err := searchingproductsv1.NewSearchProductsWithValidation(
		c.Products[0].Name,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	mapper.ClearMappings()

	res, err := c.handler.Handle(c.Ctx, query)
	c.Require().Error(err)
	c.True(customErrors.IsApplicationError(err, http.StatusInternalServerError))
	c.Nil(res)
}
