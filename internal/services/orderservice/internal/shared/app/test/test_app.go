// Package test contains the test app.
package test

import (
	"context"
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/store"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/grpc"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/bus"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/testcontainer/rabbitmq"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/testcontainer/redis"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"

	kdb "github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
	config4 "github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstoredb/config"
	config3 "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/config"
	config2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
	kurrentdb "github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/testcontainer/eventstoredb"
	mongo2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/testcontainer/mongo"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/config"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/contracts/repositories"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/aggregate"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/configurations/orders"
	ordersService "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/grpc/genproto"
)

// OrderTestApp is a struct that contains the test app.
type OrderTestApp struct{}

// OrderTestAppResult is a struct that contains the test app result.
type OrderTestAppResult struct {
	Cfg                  *config.Config
	Bus                  bus.RabbitmqBus
	Container            contracts.Container
	Logger               logger.Logger
	RabbitmqOptions      *config2.RabbitmqOptions
	EchoHTTPOptions      *config3.EchoHTTPOptions
	EventStoreDbOptions  *config4.EventStoreDbOptions
	OrderMongoRepository repositories.OrderMongoRepository
	OrderAggregateStore  store.AggregateStore[*aggregate.Order]
	OrdersServiceClient  ordersService.OrdersServiceClient
	MongoClient          *mongo.Client
	EsdbClient           *kdb.Client
	MongoDbOptions       *mongodb.MongoDbOptions
	GrpcClient           grpc.GrpcClient
}

// NewOrderTestApp is a constructor for the OrderTestApp.
func NewOrderTestApp() *OrderTestApp {
	return &OrderTestApp{}
}

// Run is a method that runs the test app.
func (a *OrderTestApp) Run(t *testing.T) (result *OrderTestAppResult) {
	t.Helper()

	lifetimeCtx := context.Background()

	// ref: https://github.com/uber-go/fx/blob/master/app_test.go
	appBuilder := NewOrderServiceTestApplicationBuilder(t)
	appBuilder.ProvideModule(orders.OrderServiceModule())

	appBuilder.Decorate(rabbitmq.RabbitmqContainerOptionsDecorator(t, lifetimeCtx))
	appBuilder.Decorate(kurrentdb.EventstoreDBContainerOptionsDecorator(t, lifetimeCtx))
	appBuilder.Decorate(mongo2.MongoContainerOptionsDecorator(t, lifetimeCtx))
	appBuilder.Decorate(redis.RedisContainerOptionsDecorator(t, lifetimeCtx))

	testApp := appBuilder.Build()

	testApp.ConfigureOrders()

	testApp.MapOrdersEndpoints()

	testApp.ResolveFunc(
		func(
			cfg *config.Config,
			bus bus.RabbitmqBus,
			logger logger.Logger,
			rabbitmqOptions *config2.RabbitmqOptions,
			echoOptions *config3.EchoHTTPOptions,
			grpcClient grpc.GrpcClient,
			eventStoreDbOptions *config4.EventStoreDbOptions,
			orderMongoRepository repositories.OrderMongoRepository,
			orderAggregateStore store.AggregateStore[*aggregate.Order],
			mongoClient *mongo.Client,
			esdbClient *kdb.Client,
			mongoDbOptions *mongodb.MongoDbOptions,
		) {
			result = &OrderTestAppResult{
				Bus:                  bus,
				Cfg:                  cfg,
				Container:            testApp,
				Logger:               logger,
				RabbitmqOptions:      rabbitmqOptions,
				MongoClient:          mongoClient,
				MongoDbOptions:       mongoDbOptions,
				EchoHTTPOptions:      echoOptions,
				EsdbClient:           esdbClient,
				EventStoreDbOptions:  eventStoreDbOptions,
				OrderMongoRepository: orderMongoRepository,
				OrderAggregateStore:  orderAggregateStore,
				OrdersServiceClient: ordersService.NewOrdersServiceClient(
					grpcClient.GetGrpcConnection(),
				),
				GrpcClient: grpcClient,
			}
		},
	)

	// we need a longer timout for up and running our testcontainers
	duration := time.Second * 300

	// short timeout for handling start hooks and setup dependencies
	startCtx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	err := testApp.Start(startCtx)
	if err != nil {
		t.Errorf("Error starting, err: %v", err)
	}

	t.Cleanup(func() {
		// short timeout for handling stop hooks
		stopCtx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		err = testApp.Stop(stopCtx)
		require.NoError(t, err)
	})

	return result
}
