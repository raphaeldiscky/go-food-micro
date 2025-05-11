package repositories

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/contracts/repositories"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/read_models"

	"github.com/elastic/go-elasticsearch/v8"
	uuid "github.com/satori/go.uuid"
)

type elasticOrderReadRepository struct {
	log           logger.Logger
	elasticClient *elasticsearch.Client
	tracer        tracing.AppTracer
}

func NewElasticOrderReadRepository(
	log logger.Logger,
	elasticClient *elasticsearch.Client,
	tracer tracing.AppTracer,
) repositories.OrderElasticRepository {
	return &elasticOrderReadRepository{log: log, elasticClient: elasticClient, tracer: tracer}
}

func (e elasticOrderReadRepository) GetAllOrders(
	ctx context.Context,
	listQuery *utils.ListQuery,
) (*utils.ListResult[*read_models.OrderReadModel], error) {
	// TODO implement me
	panic("implement me")
}

func (e elasticOrderReadRepository) SearchOrders(
	ctx context.Context,
	searchText string,
	listQuery *utils.ListQuery,
) (*utils.ListResult[*read_models.OrderReadModel], error) {
	// TODO implement me
	panic("implement me")
}

func (e elasticOrderReadRepository) GetOrderById(
	ctx context.Context,
	uuid uuid.UUID,
) (*read_models.OrderReadModel, error) {
	// TODO implement me
	panic("implement me")
}

func (e elasticOrderReadRepository) GetOrderByOrderId(
	ctx context.Context,
	uuid uuid.UUID,
) (*read_models.OrderReadModel, error) {
	// TODO implement me
	panic("implement me")
}

func (e elasticOrderReadRepository) CreateOrder(
	ctx context.Context,
	order *read_models.OrderReadModel,
) (*read_models.OrderReadModel, error) {
	// TODO implement me
	panic("implement me")
}

func (e elasticOrderReadRepository) UpdateOrder(
	ctx context.Context,
	order *read_models.OrderReadModel,
) (*read_models.OrderReadModel, error) {
	// TODO implement me
	panic("implement me")
}

func (e elasticOrderReadRepository) DeleteOrderByID(ctx context.Context, uuid uuid.UUID) error {
	// TODO implement me
	panic("implement me")
}
