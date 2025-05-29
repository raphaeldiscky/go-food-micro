//go:build unit
// +build unit

package v1

import (
	"fmt"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"
	"github.com/stretchr/testify/suite"

	gofakeit "github.com/brianvoe/gofakeit/v6"

	createProductCommand "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"
)

type createProductUnitTests struct {
	*unittest.CatalogWriteUnitTestSharedFixture
}

func TestCreateProductUnit(t *testing.T) {
	suite.Run(
		t,
		&createProductUnitTests{
			CatalogWriteUnitTestSharedFixture: unittest.NewCatalogWriteUnitTestSharedFixture(t),
		},
	)
}

// TestNewCreateProductShouldReturnNoErrorForValidInput tests the new create product should return no error for valid input.
func (c *createProductUnitTests) TestNewCreateProductShouldReturnNoErrorForValidInput() {
	name := gofakeit.Name()
	description := gofakeit.EmojiDescription()
	price := gofakeit.Price(150, 6000)

	createProduct, err := createProductCommand.NewCreateProductWithValidation(
		name,
		description,
		price,
	)
	var g interface{} = createProduct
	d, ok := g.(cqrs.Command)
	if ok {
		fmt.Println(d)
	}

	c.Assert().NotNil(createProduct)
	c.Assert().Equal(name, createProduct.Name)
	c.Assert().Equal(price, createProduct.Price)

	c.Require().NoError(err)
}

// TestNewCreateProductShouldReturnErrorForInvalidPrice tests the new create product should return error for invalid price.
func (c *createProductUnitTests) TestNewCreateProductShouldReturnErrorForInvalidPrice() {
	command, err := createProductCommand.NewCreateProductWithValidation(
		gofakeit.Name(),
		gofakeit.EmojiDescription(),
		0,
	)

	c.Require().Error(err)
	c.NotNil(command)
	c.Equal(0.0, command.Price)
}

// TestNewCreateProductShouldReturnErrorForEmptyName tests the new create product should return error for empty name.
func (c *createProductUnitTests) TestNewCreateProductShouldReturnErrorForEmptyName() {
	command, err := createProductCommand.NewCreateProductWithValidation(
		"",
		gofakeit.EmojiDescription(),
		120,
	)

	c.Require().Error(err)
	c.NotNil(command)
	c.Empty(command.Name)
}

// TestNewCreateProductShouldReturnErrorForEmptyDescription tests the new create product should return error for empty description.
func (c *createProductUnitTests) TestNewCreateProductShouldReturnErrorForEmptyDescription() {
	command, err := createProductCommand.NewCreateProductWithValidation(
		gofakeit.Name(),
		"",
		120,
	)

	c.Require().Error(err)
	c.NotNil(command)
	c.Empty(command.Description)
}
