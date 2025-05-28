// Package mappings contains the products mappings.
package mappings

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/dto"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/models"
)

// ConfigureProductsMappings is a function that configures the products mappings.
func ConfigureProductsMappings() error {
	err := mapper.CreateMap[*models.Product, *dto.ProductDto]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*models.Product, *models.Product]()
	if err != nil {
		return err
	}

	return nil
}
