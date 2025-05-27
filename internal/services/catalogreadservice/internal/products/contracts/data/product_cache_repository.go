package data

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/models"
)

type ProductCacheRepository interface {
	PutProduct(ctx context.Context, key string, product *models.Product) error
	GetProductByID(ctx context.Context, key string) (*models.Product, error)
	DeleteProduct(ctx context.Context, key string) error
	DeleteAllProducts(ctx context.Context) error
}
