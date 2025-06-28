// Package integration contains the integration test fixture.
package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/bus"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/utils"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	config2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/config"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/contracts/data"
	// Import the external event types that need to be registered
	createProductExternalEventV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/creatingproduct/v1/events/integrationevents/externalevents"
	deleteProductExternalEventV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/deletingproducts/v1/events/integrationevents/externalevents"
	updateProductExternalEventsV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/updatingproducts/v1/events/integrationevents/externalevents"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/models"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/app/test"
)

// CatalogReadIntegrationTestSharedFixture is a shared fixture for integration tests.
type CatalogReadIntegrationTestSharedFixture struct {
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

// NewCatalogReadIntegrationTestSharedFixture creates a new CatalogReadIntegrationTestSharedFixture.
func NewCatalogReadIntegrationTestSharedFixture(
	t *testing.T,
) *CatalogReadIntegrationTestSharedFixture {
	t.Helper()
	result := test.NewCatalogReadTestApp().Run(t)

	// Register message types for tests - this is crucial for message casting to work
	// Create message instances
	productCreatedMsg := &createProductExternalEventV1.ProductCreatedV1{}
	productDeletedMsg := &deleteProductExternalEventV1.ProductDeletedV1{}
	productUpdatedMsg := &updateProductExternalEventsV1.ProductUpdatedV1{}

	// Register message types using the standard utility function
	messageTypesMap := map[string]types.IMessage{
		productCreatedMsg.GetMessageTypeName(): productCreatedMsg,
		productDeletedMsg.GetMessageTypeName(): productDeletedMsg,
		productUpdatedMsg.GetMessageTypeName(): productUpdatedMsg,
	}

	utils.RegisterCustomMessageTypesToRegistry(messageTypesMap)

	result.Logger.Infow("Registered message types for integration tests", logger.Fields{
		"productCreated": productCreatedMsg.GetMessageTypeName(),
		"productDeleted": productDeletedMsg.GetMessageTypeName(),
		"productUpdated": productUpdatedMsg.GetMessageTypeName(),
	})

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
		result.RabbitmqOptions.RabbitmqHostOptions.HTTPEndPoint(),
		result.RabbitmqOptions.RabbitmqHostOptions.UserName,
		result.RabbitmqOptions.RabbitmqHostOptions.Password)
	if err != nil {
		result.Logger.Error(errors.WrapIf(err, "error in creating rabbithole client"))
		t.Fatalf("Failed to create RabbitMQ management client: %v", err)
	}

	shared := &CatalogReadIntegrationTestSharedFixture{
		Log:                    result.Logger,
		Container:              result.Container,
		Cfg:                    result.Cfg,
		RabbitmqCleaner:        rmqc,
		ProductRepository:      result.ProductRepository,
		ProductCacheRepository: result.ProductCacheRepository,
		Bus:                    result.Bus,
		rabbitmqOptions:        result.RabbitmqOptions,
		MongoOptions:           result.MongoDbOptions,
		BaseAddress:            result.EchoHTTPOptions.BasePathAddress(),
		mongoClient:            result.MongoClient,
		Tracer:                 result.Tracer,
	}

	// Start the bus with a context that has a longer timeout
	startCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to start the bus with retries and proper wait times
	var startErr error
	for i := 0; i < 5; i++ { // Increased retries
		// Use background context for the actual bus start so consumers don't get cancelled
		startErr = shared.Bus.Start(context.Background())
		if startErr == nil {
			// Success! Break out of retry loop
			result.Logger.Info("RabbitMQ bus started successfully")
			break
		}

		result.Logger.Warn(
			"Failed to start RabbitMQ bus, retrying...",
			logger.Fields{
				"attempt": i + 1,
				"error":   startErr,
			},
		)

		// Wait before retrying, but check if we should stop due to timeout
		select {
		case <-startCtx.Done():
			startErr = startCtx.Err()
			break
		case <-time.After(time.Second * time.Duration(i+1)):
			// Continue to next retry
		}
	}

	if startErr != nil {
		result.Logger.Error(errors.WrapIf(startErr, "error starting RabbitMQ bus"))
		t.Fatalf("Failed to start RabbitMQ bus after retries: %v", startErr)
	}

	return shared
}

// SetupTest sets up the test data.
func (i *CatalogReadIntegrationTestSharedFixture) SetupTest(t *testing.T) {
	t.Helper()

	i.Log.Info("SetupTest started")

	// Clean up any existing data first
	if err := i.cleanupMongoData(); err != nil {
		i.Log.Error(errors.WrapIf(err, "error in cleanup mongodb data"))
	}

	// seed data in each test
	res, err := seedData(i.mongoClient, i.MongoOptions.Database)
	if err != nil {
		i.Log.Error(errors.WrapIf(err, "error in seeding mongodb data"))
		t.Fatalf("Failed to seed test data: %v", err)
	}

	if len(res) == 0 {
		i.Log.Error("No test data was seeded")
		t.Fatalf("No test data was seeded")
	}

	i.Items = res
	i.Log.Infow(
		"Test data setup complete",
		logger.Fields{
			"product": res[0],
		},
	)
}

// TearDownTest tears down the test data.
func (i *CatalogReadIntegrationTestSharedFixture) TearDownTest() {
	i.Log.Info("TearDownTest started")

	// cleanup test containers with their hooks first
	if err := i.cleanupRabbitmqData(); err != nil {
		i.Log.Errorw(
			"Error cleaning up RabbitMQ data",
			logger.Fields{"error": err},
		)
	}

	// Stop the bus gracefully after cleanup
	if err := i.Bus.Stop(); err != nil {
		i.Log.Errorw(
			"Error stopping bus",
			logger.Fields{"error": err},
		)
	}

	// Give time for cleanup to complete
	time.Sleep(1 * time.Second)
}

// seedData seeds the test data.
func seedData(
	db *mongo.Client,
	databaseName string,
) ([]*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create 2 products for testing
	products := []*models.Product{
		{
			ID:          uuid.NewV4().String(),
			ProductID:   uuid.NewV4().String(),
			Name:        gofakeit.Name(),
			CreatedAt:   time.Now(),
			Description: gofakeit.AdjectiveDescriptive(),
			Price:       gofakeit.Price(100, 1000),
		},
		{
			ID:          uuid.NewV4().String(),
			ProductID:   uuid.NewV4().String(),
			Name:        gofakeit.Name(),
			CreatedAt:   time.Now(),
			Description: gofakeit.AdjectiveDescriptive(),
			Price:       gofakeit.Price(100, 1000),
		},
	}

	collection := db.Database(databaseName).Collection("products")

	// Insert both products
	for _, product := range products {
		_, err := collection.InsertOne(ctx, product)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to insert test product")
		}
	}

	// Retrieve all products for debugging
	var allProducts []*models.Product
	cursor, err := collection.Find(ctx, map[string]interface{}{})
	if err != nil {
		return nil, errors.WrapIf(err, "failed to find products after seeding")
	}

	if err := cursor.All(ctx, &allProducts); err != nil {
		if err := cursor.Close(ctx); err != nil {
			return nil, errors.WrapIf(err, "failed to close cursor")
		}

		return nil, errors.WrapIf(err, "failed to decode products after seeding")
	}

	err = cursor.Close(ctx)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to close cursor")
	}

	if len(allProducts) == 0 {
		return nil, errors.New("no products found after seeding")
	}

	return allProducts, nil
}

// cleanupRabbitmqData cleans up the rabbitmq data.
func (i *CatalogReadIntegrationTestSharedFixture) cleanupRabbitmqData() error {
	// https://github.com/michaelklishin/rabbit-hole
	// Get all queues
	queues, err := i.RabbitmqCleaner.ListQueuesIn(
		i.rabbitmqOptions.RabbitmqHostOptions.VirtualHost,
	)
	if err != nil {
		return errors.WrapIf(err, "failed to list queues")
	}

	// clear each queue using index-based iteration
	for idx := range queues {
		_, err = i.RabbitmqCleaner.PurgeQueue(
			i.rabbitmqOptions.RabbitmqHostOptions.VirtualHost,
			queues[idx].Name,
		)
		if err != nil {
			return errors.WrapIf(err, fmt.Sprintf("failed to purge queue: %s", queues[idx].Name))
		}
	}

	return nil
}

// cleanupMongoData cleans up the mongodb data.
func (i *CatalogReadIntegrationTestSharedFixture) cleanupMongoData() error {
	collections := []string{"products"}
	err := cleanupCollections(
		i.mongoClient,
		collections,
		i.MongoOptions.Database,
	)
	if err != nil {
		return errors.WrapIf(err, "failed to cleanup MongoDB collections")
	}

	return nil
}

// cleanupCollections cleans up the collections.
func cleanupCollections(
	db *mongo.Client,
	collections []string,
	databaseName string,
) error {
	database := db.Database(databaseName)
	ctx := context.Background()

	// Iterate over the collections and delete all documents
	for _, collectionName := range collections {
		collection := database.Collection(collectionName)
		// Use an empty filter instead of nil
		_, err := collection.DeleteMany(ctx, map[string]interface{}{})
		if err != nil {
			return errors.WrapIf(err, "failed to delete documents from collection: "+collectionName)
		}
	}

	return nil
}
