// Package mappings contains the products mappings.
package mappings

import (
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"google.golang.org/protobuf/types/known/timestamppb"

	uuid "github.com/satori/go.uuid"

	datamodel "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/data/datamodels"
	dtoV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/models"
	productsService "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/grpc/genproto"
)

// ConfigureProductsMappings is a function that configures the products mappings.
func ConfigureProductsMappings() error {
	// Time mappings
	err := mapper.CreateCustomMap[time.Time, *timestamppb.Timestamp](
		func(t time.Time) (*timestamppb.Timestamp, error) {
			return timestamppb.New(t), nil
		},
	)
	if err != nil {
		return err
	}

	err = mapper.CreateCustomMap[*timestamppb.Timestamp, time.Time](
		func(t *timestamppb.Timestamp) (time.Time, error) {
			if t == nil {
				return time.Time{}, nil
			}
			return t.AsTime(), nil
		},
	)
	if err != nil {
		return err
	}

	// Model mappings
	err = mapper.CreateMap[*models.Product, *dtoV1.ProductDto]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*dtoV1.ProductDto, *models.Product]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*datamodel.ProductDataModel, *models.Product]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*models.Product, *datamodel.ProductDataModel]()
	if err != nil {
		return err
	}

	// gRPC service mappings
	err = mapper.CreateCustomMap[*models.Product, *productsService.Product](
		func(product *models.Product) (*productsService.Product, error) {
			if product == nil {
				return nil, nil
			}

			return &productsService.Product{
				ProductID:   product.ID.String(),
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				CreatedAt:   timestamppb.New(product.CreatedAt),
				UpdatedAt:   timestamppb.New(product.UpdatedAt),
			}, nil
		},
	)
	if err != nil {
		return err
	}

	err = mapper.CreateCustomMap[*productsService.Product, *dtoV1.ProductDto](
		func(product *productsService.Product) (*dtoV1.ProductDto, error) {
			if product == nil {
				return nil, nil
			}

			id, err := uuid.FromString(product.ProductID)
			if err != nil {
				return nil, err
			}

			return &dtoV1.ProductDto{
				ID:          id,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				CreatedAt:   product.CreatedAt.AsTime(),
				UpdatedAt:   product.UpdatedAt.AsTime(),
			}, nil
		},
	)
	if err != nil {
		return err
	}

	// Add mapping from dtoV1.ProductDto to productsService.Product
	err = mapper.CreateCustomMap[*dtoV1.ProductDto, *productsService.Product](
		func(product *dtoV1.ProductDto) (*productsService.Product, error) {
			if product == nil {
				return nil, nil
			}

			return &productsService.Product{
				ProductID:   product.ID.String(),
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				CreatedAt:   timestamppb.New(product.CreatedAt),
				UpdatedAt:   timestamppb.New(product.UpdatedAt),
			}, nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}
