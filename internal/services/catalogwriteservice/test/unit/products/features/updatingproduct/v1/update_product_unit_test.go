//go:build unit
// +build unit

package v1

import (
	"testing"

	"github.com/stretchr/testify/suite"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"

	v1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/updatingproduct/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"
)

type updateProductUnitTests struct {
	*unittest.CatalogWriteUnitTestSharedFixture
}

func TestUpdateProductUnit(t *testing.T) {
	suite.Run(
		t,
		&updateProductUnitTests{
			CatalogWriteUnitTestSharedFixture: unittest.NewCatalogWriteUnitTestSharedFixture(t),
		},
	)
}

// TestNewUpdateProductShouldReturnNoErrorForValidInput tests the new update product should return no error for valid input.
func (c *updateProductUnitTests) TestNewUpdateProductShouldReturnNoErrorForValidInput() {
	id := uuid.NewV4()
	name := gofakeit.Name()
	description := gofakeit.EmojiDescription()
	price := gofakeit.Price(150, 6000)

	updateProduct, err := v1.NewUpdateProductWithValidation(id, name, description, price)

	c.Assert().NotNil(updateProduct)
	c.Assert().Equal(id, updateProduct.ProductID)
	c.Assert().Equal(name, updateProduct.Name)
	c.Assert().Equal(price, updateProduct.Price)

	c.Require().NoError(err)
}

// TestNewUpdateProductShouldReturnErrorForInvalidPrice tests the new update product should return error for invalid price.
func (c *updateProductUnitTests) TestNewUpdateProductShouldReturnErrorForInvalidPrice() {
	command, err := v1.NewUpdateProductWithValidation(
		uuid.NewV4(),
		gofakeit.Name(),
		gofakeit.EmojiDescription(),
		0,
	)

	c.Require().Error(err)
	c.NotNil(command)
	c.Equal(0.0, command.Price)
}

// TestNewUpdateProductShouldReturnErrorForEmptyName tests the new update product should return error for empty name.
func (c *updateProductUnitTests) TestNewUpdateProductShouldReturnErrorForEmptyName() {
	command, err := v1.NewUpdateProductWithValidation(
		uuid.NewV4(),
		"",
		gofakeit.EmojiDescription(),
		120,
	)

	c.Require().Error(err)
	c.NotNil(command)
	c.Empty(command.Name)
}

// TestNewUpdateProductShouldReturnErrorForEmptyDescription tests the new update product should return error for empty description.
func (c *updateProductUnitTests) TestNewUpdateProductShouldReturnErrorForEmptyDescription() {
	command, err := v1.NewUpdateProductWithValidation(uuid.NewV4(), gofakeit.Name(), "", 120)

	c.Require().Error(err)
	c.NotNil(command)
	c.Empty(command.Description)
}
