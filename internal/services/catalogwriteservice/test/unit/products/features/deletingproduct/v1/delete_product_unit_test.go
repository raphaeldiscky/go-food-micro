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

func (c *deleteProductUnitTests) Test_New_Delete_Product_Should_Return_No_Error_For_Valid_Input() {
	id := uuid.NewV4()

	command, err := v1.NewDeleteProductWithValidation(id)

	c.Assert().NotNil(command)
	c.Assert().Equal(command.ProductID, id)
	c.Require().NoError(err)
}

func (c *deleteProductUnitTests) Test_New_Delete_Product_Should_Return_Error_For_Invalid_Id() {
	command, err := v1.NewDeleteProductWithValidation(uuid.Nil)

	c.Require().Error(err)
	c.NotNil(command)
	c.Equal(uuid.Nil, command.ProductID)
}
