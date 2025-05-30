package data

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/models"
)

// ProductRepository is a contract for the product repository.
type ProductRepository interface {
	GetAllProducts(
		ctx context.Context,
		listQuery *utils.ListQuery,
	) (*utils.ListResult[*models.Product], error)
	SearchProducts(
		ctx context.Context,
		searchText string,
		listQuery *utils.ListQuery,
	) (*utils.ListResult[*models.Product], error)
	GetProductByID(ctx context.Context, uuid string) (*models.Product, error)
	GetProductByProductID(ctx context.Context, uuid string) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	DeleteProductByID(ctx context.Context, uuid string) error
}
