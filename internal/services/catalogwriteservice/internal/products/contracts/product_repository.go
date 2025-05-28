// Package contracts contains the product repository contract.
package contracts

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/models"
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
	GetProductByID(ctx context.Context, uuid uuid.UUID) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	DeleteProductByID(ctx context.Context, uuid uuid.UUID) error
}
