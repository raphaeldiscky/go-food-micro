package repositories

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/readmodels"
)

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
	GetOrderByOrderId(ctx context.Context, orderId uuid.UUID) (*readmodels.OrderReadModel, error)
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

type OrderElasticRepository interface {
	orderReadRepository
}

type OrderMongoRepository interface {
	orderReadRepository
}
