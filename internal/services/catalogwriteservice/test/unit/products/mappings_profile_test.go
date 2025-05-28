//go:build unit
// +build unit

package products

import (
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	dtoV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/models"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/unittest"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type mappingProfileUnitTests struct {
	*unittest.CatalogWriteUnitTestSharedFixture
}

func TestMappingProfileUnit(t *testing.T) {
	suite.Run(
		t,
		&mappingProfileUnitTests{CatalogWriteUnitTestSharedFixture: unittest.NewCatalogWriteUnitTestSharedFixture(t)},
	)
}

func (m *mappingProfileUnitTests) Test_Mappings() {
	productModel := &models.Product{
		ID:          uuid.NewV4(),
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	productDto := &dtoV1.ProductDto{
		ID:          uuid.NewV4(),
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	m.Run("Should_Map_Product_To_ProductDto", func() {
		d, err := mapper.Map[*dtoV1.ProductDto](productModel)
		m.Require().NoError(err)
		m.Equal(productModel.ID, d.ID)
		m.Equal(productModel.Name, d.Name)
	})

	m.Run("Should_Map_Nil_Product_To_ProductDto", func() {
		d, err := mapper.Map[*dtoV1.ProductDto](*new(models.Product))
		m.Require().NoError(err)
		m.Nil(d)
	})

	m.Run("Should_Map_ProductDto_To_Product", func() {
		d, err := mapper.Map[*models.Product](productDto)
		m.Require().NoError(err)
		m.Equal(productDto.ID, d.ID)
		m.Equal(productDto.Name, d.Name)
	})

	m.Run("Should_Map_Nil_ProductDto_To_Product", func() {
		d, err := mapper.Map[*models.Product](*new(dtoV1.ProductDto))
		m.Require().NoError(err)
		m.Nil(d)
	})
}
