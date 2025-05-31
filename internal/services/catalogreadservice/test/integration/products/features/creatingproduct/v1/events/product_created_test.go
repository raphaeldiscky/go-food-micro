//go:build integration
// +build integration

package events

// https://github.com/smartystreets/goconvey/wiki

import (
	"context"
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/messaging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"

	externalEvents "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/creatingproduct/v1/events/integrationevents/externalevents"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/models"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/testfixture/integration"
)

func TestProductCreatedConsumer(t *testing.T) {
	// Initialize the shared fixture for integration tests
	integrationTestSharedFixture := integration.NewCatalogReadIntegrationTestSharedFixture(t)
	require.NotNil(
		t,
		integrationTestSharedFixture,
		"Integration test shared fixture should not be nil",
	)

	// The bus is already started in NewCatalogReadIntegrationTestSharedFixture
	// No need to start it again here

	// Wait for the bus to be ready and consumers to be initialized
	time.Sleep(5 * time.Second)

	t.Run("should consume ProductCreated event and create product in MongoDB", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		fakeProduct := &externalEvents.ProductCreatedV1{
			Message: &types.Message{
				MessageId: uuid.NewV4().String(),
				Created:   time.Now().UTC(),
			},
			ProductID:   uuid.NewV4().String(),
			Name:        gofakeit.Name(),
			Description: gofakeit.Sentence(10),
			Price:       gofakeit.Float64Range(1, 1000),
			CreatedAt:   time.Now().UTC(),
		}

		// Log the test setup
		integrationTestSharedFixture.Log.Infow(
			"Starting ProductCreated consumer test",
			logger.Fields{
				"productId": fakeProduct.ProductID,
				"name":      fakeProduct.Name,
				"price":     fakeProduct.Price,
			},
		)

		// Act - Publish message with retries
		var publishErr error
		for i := 0; i < 3; i++ {
			publishErr = integrationTestSharedFixture.Bus.PublishMessage(ctx, fakeProduct, nil)
			if publishErr == nil {
				break
			}
			integrationTestSharedFixture.Log.Warnw(
				"Failed to publish message, retrying...",
				logger.Fields{
					"attempt": i + 1,
					"error":   publishErr,
				},
			)
			time.Sleep(time.Second)
		}
		require.NoError(t, publishErr, "Failed to publish message after retries")

		// Assert - Wait for the message to be consumed with retries
		hypothesis := messaging.ShouldConsume(
			ctx,
			integrationTestSharedFixture.Bus,
			func(msg *externalEvents.ProductCreatedV1) bool {
				if msg == nil {
					integrationTestSharedFixture.Log.Error("Received nil message")
					return false
				}
				integrationTestSharedFixture.Log.Infow(
					"Received message",
					logger.Fields{
						"productId": msg.ProductID,
						"name":      msg.Name,
						"price":     msg.Price,
					},
				)
				return msg.ProductID == fakeProduct.ProductID
			},
		)

		// Validate the hypothesis with a timeout
		err := hypothesis.Validate(ctx, "Message should be consumed", 30*time.Second)
		require.NoError(t, err, "Message was not consumed within timeout")

		// Wait for the product to be created in the database with retries
		var product *models.Product
		var dbErr error
		for i := 0; i < 10; i++ {
			product, dbErr = integrationTestSharedFixture.ProductRepository.GetProductByProductID(
				ctx,
				fakeProduct.ProductID,
			)
			if dbErr == nil && product != nil {
				break
			}
			integrationTestSharedFixture.Log.Warnw(
				"Product not found in database, retrying...",
				logger.Fields{
					"attempt": i + 1,
					"error":   dbErr,
				},
			)
			time.Sleep(time.Second)
		}
		require.NoError(t, dbErr, "Failed to get product from database after retries")
		require.NotNil(t, product, "Product should be stored in database")

		// Verify the product data
		assert.Equal(t, fakeProduct.ProductID, product.ProductID)
		assert.Equal(t, fakeProduct.Name, product.Name)
		assert.Equal(t, fakeProduct.Description, product.Description)
		assert.Equal(t, fakeProduct.Price, product.Price)
		assert.Equal(t, fakeProduct.CreatedAt.Unix(), product.CreatedAt.Unix())

		// Verify no duplicate products
		retrievedProduct, err := integrationTestSharedFixture.ProductRepository.GetProductByProductID(
			ctx,
			fakeProduct.ProductID,
		)
		require.NoError(t, err, "Failed to get product from database")
		require.NotNil(t, retrievedProduct, "Product should exist in database")
		assert.Equal(
			t,
			fakeProduct.ProductID,
			retrievedProduct.ProductID,
			"Should have exactly one product with the given ID",
		)

		integrationTestSharedFixture.Log.Info("ProductCreated consumer test completed successfully")
	})
}
