package integration

import (
	"context"
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/bus"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
	config2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/config"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/contracts/data"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/models"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/app/test"

	"emperror.dev/errors"
	"github.com/brianvoe/gofakeit/v6"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/trace"
)

type IntegrationTestSharedFixture struct {
	Cfg                    *config.Config
	Log                    logger.Logger
	Bus                    bus.Bus
	ProductRepository      data.ProductRepository
	ProductCacheRepository data.ProductCacheRepository
	Container              contracts.Container
	RabbitmqCleaner        *rabbithole.Client
	rabbitmqOptions        *config2.RabbitmqOptions
	MongoOptions           *mongodb.MongoDbOptions
	BaseAddress            string
	mongoClient            *mongo.Client
	Items                  []*models.Product
	Tracer                 trace.Tracer
}

func NewIntegrationTestSharedFixture(
	t *testing.T,
) *IntegrationTestSharedFixture {
	result := test.NewTestApp().Run(t)

	// Log the RabbitMQ connection details for debugging
	result.Logger.Infow(
		"RabbitMQ connection details",
		logger.Fields{
			"host":     result.RabbitmqOptions.RabbitmqHostOptions.HostName,
			"port":     result.RabbitmqOptions.RabbitmqHostOptions.Port,
			"httpPort": result.RabbitmqOptions.RabbitmqHostOptions.HttpPort,
			"user":     result.RabbitmqOptions.RabbitmqHostOptions.UserName,
			"vhost":    result.RabbitmqOptions.RabbitmqHostOptions.VirtualHost,
		},
	)

	// https://github.com/michaelklishin/rabbit-hole
	rmqc, err := rabbithole.NewClient(
		result.RabbitmqOptions.RabbitmqHostOptions.HttpEndPoint(),
		result.RabbitmqOptions.RabbitmqHostOptions.UserName,
		result.RabbitmqOptions.RabbitmqHostOptions.Password)
	if err != nil {
		result.Logger.Error(errors.WrapIf(err, "error in creating rabbithole client"))
		t.Fatalf("Failed to create RabbitMQ management client: %v", err)
	}

	shared := &IntegrationTestSharedFixture{
		Log:                    result.Logger,
		Container:              result.Container,
		Cfg:                    result.Cfg,
		RabbitmqCleaner:        rmqc,
		ProductRepository:      result.ProductRepository,
		ProductCacheRepository: result.ProductCacheRepository,
		Bus:                    result.Bus,
		rabbitmqOptions:        result.RabbitmqOptions,
		MongoOptions:           result.MongoDbOptions,
		BaseAddress:            result.EchoHttpOptions.BasePathAddress(),
		mongoClient:            result.MongoClient,
		Tracer:                 result.Tracer,
	}

	// Start the bus with a context that has a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to start the bus with retries
	var startErr error
	for i := 0; i < 3; i++ {
		startErr = shared.Bus.Start(ctx)
		if startErr == nil {
			break
		}
		result.Logger.Warn(
			"Failed to start RabbitMQ bus, retrying...",
			logger.Fields{
				"attempt": i + 1,
				"error":   startErr,
			},
		)
		time.Sleep(2 * time.Second)
	}

	if startErr != nil {
		result.Logger.Error(errors.WrapIf(startErr, "error starting RabbitMQ bus"))
		t.Fatalf("Failed to start RabbitMQ bus after retries: %v", startErr)
	}

	// Wait for the bus to be ready
	time.Sleep(5 * time.Second)

	return shared
}

func (i *IntegrationTestSharedFixture) SetupTest() {
	i.Log.Info("SetupTest started")

	// seed data in each test
	res, err := seedData(i.mongoClient, i.MongoOptions.Database)
	if err != nil {
		i.Log.Error(errors.WrapIf(err, "error in seeding mongodb data"))
	}

	i.Items = res
}

func (i *IntegrationTestSharedFixture) TearDownTest() {
	i.Log.Info("TearDownTest started")

	// cleanup test containers with their hooks
	if err := i.cleanupRabbitmqData(); err != nil {
		i.Log.Error(errors.WrapIf(err, "error in cleanup rabbitmq data"))
	}

	if err := i.cleanupMongoData(); err != nil {
		i.Log.Error(errors.WrapIf(err, "error in cleanup mongodb data"))
	}
}

func seedData(
	db *mongo.Client,
	databaseName string,
) ([]*models.Product, error) {
	ctx := context.Background()

	products := []*models.Product{
		{
			Id:          uuid.NewV4().String(),
			ProductId:   uuid.NewV4().String(),
			Name:        gofakeit.Name(),
			CreatedAt:   time.Now(),
			Description: gofakeit.AdjectiveDescriptive(),
			Price:       gofakeit.Price(100, 1000),
		},
		{
			Id:          uuid.NewV4().String(),
			ProductId:   uuid.NewV4().String(),
			Name:        gofakeit.Name(),
			CreatedAt:   time.Now(),
			Description: gofakeit.AdjectiveDescriptive(),
			Price:       gofakeit.Price(100, 1000),
		},
	}

	//// https://go.dev/doc/faq#convert_slice_of_interface
	productsData := make([]interface{}, len(products))

	for i, v := range products {
		productsData[i] = v
	}

	collection := db.Database(databaseName).Collection("products")
	_, err := collection.InsertMany(
		context.Background(),
		productsData,
		&options.InsertManyOptions{},
	)
	if err != nil {
		return nil, errors.WrapIf(err, "error in seed database")
	}

	result, err := mongodb.Paginate[*models.Product](
		ctx,
		utils.NewListQuery(10, 1),
		collection,
		nil,
	)

	return result.Items, nil
}

func (i *IntegrationTestSharedFixture) cleanupRabbitmqData() error {
	// https://github.com/michaelklishin/rabbit-hole
	// Get all queues
	queues, err := i.RabbitmqCleaner.ListQueuesIn(
		i.rabbitmqOptions.RabbitmqHostOptions.VirtualHost,
	)
	if err != nil {
		return err
	}

	// clear each queue
	for _, queue := range queues {
		_, err = i.RabbitmqCleaner.PurgeQueue(
			i.rabbitmqOptions.RabbitmqHostOptions.VirtualHost,
			queue.Name,
		)

		return err
	}

	return nil
}

func (i *IntegrationTestSharedFixture) cleanupMongoData() error {
	collections := []string{"products"}
	err := cleanupCollections(
		i.mongoClient,
		collections,
		i.MongoOptions.Database,
	)

	return err
}

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
