//go:build unit
// +build unit

package v1

import (
	"testing"

	"github.com/stretchr/testify/suite"

	uuid "github.com/satori/go.uuid"

	v1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/deletingproduct/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"
)

type deleteProductUnitTests struct {
	*unittest.CatalogWriteUnitTestSharedFixture
}

func TestDeleteProductByIdUnit(t *testing.T) {
	suite.Run(
		t,
		&deleteProductUnitTests{
			CatalogWriteUnitTestSharedFixture: unittest.NewCatalogWriteUnitTestSharedFixture(t),
		},
	)
}

// TestNewDeleteProductShouldReturnNoErrorForValidInput tests the new delete product should return no error for valid input.
func (c *deleteProductUnitTests) TestNewDeleteProductShouldReturnNoErrorForValidInput() {
	id := uuid.NewV4()

	command, err := v1.NewDeleteProductWithValidation(id)

	c.Assert().NotNil(command)
	c.Assert().Equal(command.ProductID, id)
	c.Require().NoError(err)
}

// TestNewDeleteProductShouldReturnErrorForInvalidId tests the new delete product should return error for invalid id.
func (c *deleteProductUnitTests) TestNewDeleteProductShouldReturnErrorForInvalidId() {
	command, err := v1.NewDeleteProductWithValidation(uuid.Nil)

	c.Require().Error(err)
	c.NotNil(command)
	c.Equal(uuid.Nil, command.ProductID)
}
