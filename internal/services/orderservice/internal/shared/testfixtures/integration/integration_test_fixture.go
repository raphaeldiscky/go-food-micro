// Package integration contains the integration test fixture.
package integration

import (
	"context"
	"testing"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/bus"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/store"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	config3 "github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstoredb/config"

	config2 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/config"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/contracts/repositories"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/aggregate"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/readmodels"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/app/test"
	contracts2 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/contracts"
	ordersService "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/grpc/genproto"
)

const (
	orderCollection = "orders"
)

// OrderIntegrationTestSharedFixture is the integration test fixture.
type OrderIntegrationTestSharedFixture struct {
	OrderAggregateStore  store.AggregateStore[*aggregate.Order]
	OrderMongoRepository repositories.OrderMongoRepository
	OrdersMetrics        contracts2.OrdersMetrics
	Cfg                  *config2.Config
	Log                  logger.Logger
	Bus                  bus.Bus
	Container            contracts.Container
	RabbitmqCleaner      *rabbithole.Client
	rabbitmqOptions      *config.RabbitmqOptions
	BaseAddress          string
	mongoClient          *mongo.Client
	MongoDbOptions       *mongodb.MongoDbOptions
	EventStoreDbOptions  *config3.EventStoreDbOptions
	Items                []*readmodels.OrderReadModel
	OrdersServiceClient  ordersService.OrdersServiceClient
}

// NewIntegrationTestSharedFixture creates a new integration test fixture.
func NewIntegrationTestSharedFixture(
	t *testing.T,
) *OrderIntegrationTestSharedFixture {
	t.Helper()
	result := test.NewOrderTestApp().Run(t)

	// https://github.com/michaelklishin/rabbit-hole
	rmqc, err := rabbithole.NewClient(
		result.RabbitmqOptions.RabbitmqHostOptions.HTTPEndPoint(),
		result.RabbitmqOptions.RabbitmqHostOptions.UserName,
		result.RabbitmqOptions.RabbitmqHostOptions.Password)
	if err != nil {
		result.Logger.Error(
			errors.WrapIf(err, "error in creating rabbithole client"),
		)
	}
	shared := &OrderIntegrationTestSharedFixture{
		Log:                  result.Logger,
		Container:            result.Container,
		Cfg:                  result.Cfg,
		RabbitmqCleaner:      rmqc,
		OrderMongoRepository: result.OrderMongoRepository,
		OrderAggregateStore:  result.OrderAggregateStore,
		MongoDbOptions:       result.MongoDbOptions,
		EventStoreDbOptions:  result.EventStoreDbOptions,
		mongoClient:          result.MongoClient,
		Bus:                  result.Bus,
		rabbitmqOptions:      result.RabbitmqOptions,
		BaseAddress:          result.EchoHTTPOptions.BasePathAddress(),
		OrdersServiceClient:  result.OrdersServiceClient,
	}

	return shared
}

// SetupTest setups the test.
func (i *OrderIntegrationTestSharedFixture) SetupTest() {
	i.Log.Info("SetupTest started")

	// seed data in each test
	res, err := seedReadModelData(i.mongoClient, i.MongoDbOptions.Database)
	if err != nil {
		i.Log.Error(errors.WrapIf(err, "error in seeding mongodb data"))
	}
	i.Items = res
}

// TearDownTest tears down the test.
func (i *OrderIntegrationTestSharedFixture) TearDownTest() {
	i.Log.Info("TearDownTest started")

	// cleanup test containers with their hooks
	if err := i.cleanupRabbitmqData(); err != nil {
		i.Log.Error(errors.WrapIf(err, "error in cleanup rabbitmq data"))
	}

	if err := i.cleanupMongoData(); err != nil {
		i.Log.Error(errors.WrapIf(err, "error in cleanup mongodb data"))
	}
}

// cleanupRabbitmqData cleans up the rabbitmq data.
func (i *OrderIntegrationTestSharedFixture) cleanupRabbitmqData() error {
	// https://github.com/michaelklishin/rabbit-hole
	// Get all queues
	queues, err := i.RabbitmqCleaner.ListQueuesIn(
		i.rabbitmqOptions.RabbitmqHostOptions.VirtualHost,
	)
	if err != nil {
		return err
	}

	// clear each queue
	for idx := range queues {
		_, err = i.RabbitmqCleaner.PurgeQueue(
			i.rabbitmqOptions.RabbitmqHostOptions.VirtualHost,
			queues[idx].Name,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// cleanupMongoData cleans up the mongodb data.
func (i *OrderIntegrationTestSharedFixture) cleanupMongoData() error {
	collections := []string{orderCollection}
	err := cleanupCollections(
		i.mongoClient,
		collections,
		i.MongoDbOptions.Database,
	)

	return err
}

// cleanupCollections cleans up the collections.
func cleanupCollections(
	db *mongo.Client,
	collections []string,
	databaseName string,
) error {
	database := db.Database(databaseName)
	ctx := context.Background()

	// Iterate over the collections and delete all collections
	for _, collection := range collections {
		collection := database.Collection(collection)

		err := collection.Drop(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// seedReadModelData seeds the read model data.
func seedReadModelData(
	db *mongo.Client,
	databaseName string,
) ([]*readmodels.OrderReadModel, error) {
	ctx := context.Background()

	orders := []*readmodels.OrderReadModel{
		{
			ID:              gofakeit.UUID(),
			OrderID:         gofakeit.UUID(),
			ShopItems:       generateShopItems(),
			AccountEmail:    gofakeit.Email(),
			DeliveryAddress: gofakeit.Address().Address,
			CancelReason:    gofakeit.Sentence(5),
			TotalPrice:      gofakeit.Float64Range(10, 100),
			DeliveredTime:   gofakeit.Date(),
			Paid:            gofakeit.Bool(),
			Submitted:       gofakeit.Bool(),
			Completed:       gofakeit.Bool(),
			Canceled:        gofakeit.Bool(),
			PaymentID:       gofakeit.UUID(),
			CreatedAt:       gofakeit.Date(),
			UpdatedAt:       gofakeit.Date(),
		},
		{
			ID:              gofakeit.UUID(),
			OrderID:         gofakeit.UUID(),
			ShopItems:       generateShopItems(),
			AccountEmail:    gofakeit.Email(),
			DeliveryAddress: gofakeit.Address().Address,
			CancelReason:    gofakeit.Sentence(5),
			TotalPrice:      gofakeit.Float64Range(10, 100),
			DeliveredTime:   gofakeit.Date(),
			Paid:            gofakeit.Bool(),
			Submitted:       gofakeit.Bool(),
			Completed:       gofakeit.Bool(),
			Canceled:        gofakeit.Bool(),
			PaymentID:       gofakeit.UUID(),
			CreatedAt:       gofakeit.Date(),
			UpdatedAt:       gofakeit.Date(),
		},
	}

	//// https://go.dev/doc/faq#convert_slice_of_interface
	data := make([]interface{}, len(orders))
	for i, v := range orders {
		data[i] = v
	}

	collection := db.Database(databaseName).Collection("orders")
	_, err := collection.InsertMany(
		context.Background(),
		data,
		&options.InsertManyOptions{},
	)
	if err != nil {
		return nil, errors.WrapIf(err, "error in seed database")
	}

	result, err := mongodb.Paginate[*readmodels.OrderReadModel](
		ctx,
		utils.NewListQuery(10, 1),
		collection,
		nil,
	)
	if err != nil {
		return nil, errors.WrapIf(err, "error in paginate mongodb data")
	}

	return result.Items, nil
}

// generateShopItems generates the shop items.
func generateShopItems() []*readmodels.ShopItemReadModel {
	var shopItems []*readmodels.ShopItemReadModel

	for i := 0; i < 3; i++ {
		shopItem := &readmodels.ShopItemReadModel{
			Title:       gofakeit.Word(),
			Description: gofakeit.Sentence(3),
			Quantity:    uint64(gofakeit.UintRange(1, 100)),
			Price:       gofakeit.Float64Range(1, 50),
		}

		shopItems = append(shopItems, shopItem)
	}

	return shopItems
}
