//go:build unit
// +build unit

package v1

import (
	"testing"

	"github.com/stretchr/testify/suite"

	uuid "github.com/satori/go.uuid"

	getProductByIdQuery "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/gettingproductbyid/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"
)

type getProductByIdUnitTests struct {
	*unittest.CatalogWriteUnitTestSharedFixture
}

func TestGetProductByIdUnit(t *testing.T) {
	suite.Run(
		t,
		&getProductByIdUnitTests{
			CatalogWriteUnitTestSharedFixture: unittest.NewCatalogWriteUnitTestSharedFixture(t),
		},
	)
}

// TestNewGetProductByIdShouldReturnNoErrorForValidInput tests the new get product by id should return no error for valid input.
func (c *getProductByIdUnitTests) TestNewGetProductByIdShouldReturnNoErrorForValidInput() {
	id := uuid.NewV4()

	query, err := getProductByIdQuery.NewGetProductByIDWithValidation(id)

	c.Assert().NotNil(query)
	c.Assert().Equal(query.ProductID, id)
	c.Require().NoError(err)
}

// TestNewGetProductByIdShouldReturnErrorForInvalidId tests the new get product by id should return error for invalid id.
func (c *getProductByIdUnitTests) TestNewGetProductByIdShouldReturnErrorForInvalidId() {
	query, err := getProductByIdQuery.NewGetProductByIDWithValidation(uuid.Nil)

	c.Require().Error(err)
	c.NotNil(query)
	c.Equal(uuid.Nil, query.ProductID)
}
