// Package repositories contains the order repository.
package repositories

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/readmodels"
)

// orderReadRepository is the read repository for the order.
type orderReadRepository interface {
	GetAllOrders(
		ctx context.Context,
		listQuery *utils.ListQuery,
	) (*utils.ListResult[*readmodels.OrderReadModel], error)
	SearchOrders(
		ctx context.Context,
		searchText string,
		listQuery *utils.ListQuery,
	) (*utils.ListResult[*readmodels.OrderReadModel], error)
	GetOrderByID(ctx context.Context, uuid uuid.UUID) (*readmodels.OrderReadModel, error)
	GetOrderByOrderID(ctx context.Context, orderID uuid.UUID) (*readmodels.OrderReadModel, error)
	CreateOrder(
		ctx context.Context,
		order *readmodels.OrderReadModel,
	) (*readmodels.OrderReadModel, error)
	UpdateOrder(
		ctx context.Context,
		order *readmodels.OrderReadModel,
	) (*readmodels.OrderReadModel, error)
	DeleteOrderByID(ctx context.Context, uuid uuid.UUID) error
}

// OrderElasticRepository is the elastic repository for the order.
type OrderElasticRepository interface {
	orderReadRepository
}

// OrderMongoRepository is the mongo repository for the order.
type OrderMongoRepository interface {
	orderReadRepository
}
