//go:build integration
// +build integration

package events

// https://github.com/smartystreets/goconvey/wiki

import (
	"context"
	"go-food-micro/internal/pkg/core/messaging"
	"go-food-micro/internal/pkg/core/types"
	"go-food-micro/internal/services/catalogreadservice/test/integration/shared"
	"go-food-micro/internal/services/catalogreadservice/test/integration/shared/fixtures"
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/messaging"
	externalEvents "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/creating_product/v1/events/integrationevents/externalevents"

	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProductCreatedConsumer(t *testing.T) {
	// Initialize the shared fixture for integration tests
	integrationTestSharedFixture := fixtures.NewIntegrationTestSharedFixture(t)
	require.NotNil(t, integrationTestSharedFixture, "Integration test shared fixture should not be nil")

	// Start the RabbitMQ bus manually to ensure it's running before the test
	err := integrationTestSharedFixture.Bus.Start(context.Background())
	require.NoError(t, err, "Failed to start RabbitMQ bus")

	// Ensure proper cleanup
	defer func() {
		// Give some time for messages to be processed before cleanup
		time.Sleep(2 * time.Second)
		err := integrationTestSharedFixture.Bus.Stop(context.Background())
		assert.NoError(t, err, "Failed to stop RabbitMQ bus")
	}()

	// Wait for the bus to be ready
	time.Sleep(2 * time.Second)

	t.Run("should consume ProductCreated event and create product in MongoDB", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		fakeProduct := &externalEvents.ProductCreatedV1{
			Message: &types.Message{
				Id:        uuid.New().String(),
				Timestamp: time.Now().UTC(),
			},
			ProductId:   uuid.New().String(),
			Name:        gofakeit.Name(),
			Description: gofakeit.Sentence(10),
			Price:       gofakeit.Float64Range(1, 1000),
			CreatedAt:   time.Now().UTC(),
		}

		// Act
		err := integrationTestSharedFixture.Bus.PublishMessage(ctx, fakeProduct, nil)
		require.NoError(t, err, "Failed to publish message")

		// Assert - Wait for the message to be consumed with retries
		var consumed bool
		for i := 0; i < 5; i++ {
			consumed, err = messaging.ShouldConsume(ctx, integrationTestSharedFixture.Bus, fakeProduct, 5*time.Second)
			if consumed && err == nil {
				break
			}
			time.Sleep(time.Second)
		}
		assert.NoError(t, err, "Failed to consume message")
		assert.True(t, consumed, "Message should be consumed")

		// Wait for the product to be created in the database with retries
		var product *shared.Product
		for i := 0; i < 5; i++ {
			product, err = integrationTestSharedFixture.MongoDB.GetProduct(ctx, fakeProduct.ProductId)
			if err == nil && product != nil {
				break
			}
			time.Sleep(time.Second)
		}
		require.NoError(t, err, "Failed to get product from database")
		require.NotNil(t, product, "Product should be stored in database")

		// Verify the product data
		assert.Equal(t, fakeProduct.ProductId, product.Id)
		assert.Equal(t, fakeProduct.Name, product.Name)
		assert.Equal(t, fakeProduct.Description, product.Description)
		assert.Equal(t, fakeProduct.Price, product.Price)
		assert.Equal(t, fakeProduct.CreatedAt.Unix(), product.CreatedAt.Unix())

		// Verify no duplicate products
		products, err := integrationTestSharedFixture.MongoDB.GetProducts(ctx)
		require.NoError(t, err, "Failed to get products from database")
		count := 0
		for _, p := range products {
			if p.Id == fakeProduct.ProductId {
				count++
			}
		}
		assert.Equal(t, 1, count, "Should have exactly one product with the given ID")
	})
}
