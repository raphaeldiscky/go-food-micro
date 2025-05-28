//go:build unit
// +build unit

package v1

import (
	"testing"

	getProductByIdQuery "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/gettingproductbyid/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type getProductByIdUnitTests struct {
	*unittest.CatalogWriteUnitTestSharedFixture
}

func TestGetProductByIdUnit(t *testing.T) {
	suite.Run(
		t,
		&getProductByIdUnitTests{CatalogWriteUnitTestSharedFixture: unittest.NewCatalogWriteUnitTestSharedFixture(t)},
	)
}

func (c *getProductByIdUnitTests) Test_New_Get_Product_By_Id_Should_Return_No_Error_For_Valid_Input() {
	id := uuid.NewV4()

	query, err := getProductByIdQuery.NewGetProductByIDWithValidation(id)

	c.Assert().NotNil(query)
	c.Assert().Equal(query.ProductID, id)
	c.Require().NoError(err)
}

func (c *getProductByIdUnitTests) Test_New_Get_Product_By_Id_Should_Return_Error_For_Invalid_Id() {
	query, err := getProductByIdQuery.NewGetProductByIDWithValidation(uuid.Nil)

	c.Require().Error(err)
	c.NotNil(query)
	c.Equal(uuid.Nil, query.ProductID)
}
