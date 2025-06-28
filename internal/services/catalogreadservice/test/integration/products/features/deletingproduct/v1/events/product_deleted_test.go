//go:build integration
// +build integration

package events

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/messaging"

	. "github.com/smartystreets/goconvey/convey"

	testutils "github.com/raphaeldiscky/go-food-micro/internal/pkg/test/utils"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/deletingproducts/v1/events/integrationevents/externalevents"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/models"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/testfixture/integration"
)

func TestProductDeleted(t *testing.T) {
	// Setup and initialization code here.
	integrationTestSharedFixture := integration.NewCatalogReadIntegrationTestSharedFixture(t)
	// The bus is already started by the integration fixture, no need to start it again
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)

	Convey("Product Deleted Feature", t, func() {
		ctx := context.Background()
		// will execute with each subtest
		integrationTestSharedFixture.SetupTest(t)

		// First verify the product exists and log its details
		existingProduct, err := integrationTestSharedFixture.ProductRepository.GetProductByProductID(
			ctx,
			integrationTestSharedFixture.Items[0].ProductID,
		)
		So(err, ShouldBeNil)

		So(existingProduct, ShouldNotBeNil)
		integrationTestSharedFixture.Log.Infow(
			"Starting product deletion test",
			logger.Fields{
				"product":   existingProduct,
				"productID": existingProduct.ProductID,
				"name":      existingProduct.Name,
			},
		)

		Convey("Consume ProductDeleted event by consumer", func() {
			event := &externalevents.ProductDeletedV1{
				Message:   types.NewMessage(uuid.NewV4().String()),
				ProductID: integrationTestSharedFixture.Items[0].ProductID,
			}

			// Add extensive debug logging
			integrationTestSharedFixture.Log.Infow(
				"Creating hypothesis for ProductDeleted message consumption",
				logger.Fields{
					"messageType":     fmt.Sprintf("%T", event),
					"messageTypeName": event.GetMessageTypeName(),
					"productID":       event.ProductID,
				},
			)

			hypothesis := messaging.ShouldConsume[*externalevents.ProductDeletedV1](
				ctx,
				integrationTestSharedFixture.Bus,
				nil,
			)

			Convey("When a ProductDeleted event consumed", func() {
				integrationTestSharedFixture.Log.Infow(
					"About to publish ProductDeleted message",
					logger.Fields{
						"messageId":       event.GeMessageId(),
						"messageType":     fmt.Sprintf("%T", event),
						"messageTypeName": event.GetMessageTypeName(),
						"productID":       event.ProductID,
					},
				)

				err := integrationTestSharedFixture.Bus.PublishMessage(
					ctx,
					event,
					nil,
				)
				So(err, ShouldBeNil)

				integrationTestSharedFixture.Log.Infow(
					"Successfully published ProductDeleted message",
					logger.Fields{
						"messageId": event.GeMessageId(),
						"productID": event.ProductID,
						"error":     err,
					},
				)

				Convey("Then it should consume the ProductDeleted event", func() {
					integrationTestSharedFixture.Log.Infow(
						"Starting hypothesis validation with 30 second timeout",
						logger.Fields{},
					)
					hypothesis.Validate(
						ctx,
						"there is no consumed ProductDeleted message",
						30*time.Second,
					)
				})
			})
		})

		Convey("Delete product in mongo database when a ProductDeleted event consumed", func() {
			// Add small delay to ensure consumers are ready after the previous test
			time.Sleep(2 * time.Second)

			// Restart the bus to ensure consumers are reconnected properly
			integrationTestSharedFixture.Log.Info(
				"Second test: Restarting bus to ensure consumers are ready",
			)
			integrationTestSharedFixture.Bus.Stop()
			time.Sleep(1 * time.Second)

			// Start the bus again
			err := integrationTestSharedFixture.Bus.Start(context.Background())
			if err != nil {
				integrationTestSharedFixture.Log.Errorw(
					"Failed to restart bus",
					logger.Fields{"error": err},
				)
			} else {
				integrationTestSharedFixture.Log.Info("Bus restarted successfully")
			}

			// Wait for consumers to be ready
			time.Sleep(3 * time.Second)

			event := &externalevents.ProductDeletedV1{
				Message:   types.NewMessage(uuid.NewV4().String()),
				ProductID: integrationTestSharedFixture.Items[0].ProductID,
			}

			integrationTestSharedFixture.Log.Infow(
				"Second test: About to publish ProductDeleted message",
				logger.Fields{
					"messageId":       event.GeMessageId(),
					"messageType":     fmt.Sprintf("%T", event),
					"messageTypeName": event.GetMessageTypeName(),
					"productID":       event.ProductID,
				},
			)

			Convey("When a ProductDeleted event consumed", func() {
				err := integrationTestSharedFixture.Bus.PublishMessage(
					ctx,
					event,
					nil,
				)
				So(err, ShouldBeNil)

				integrationTestSharedFixture.Log.Infow(
					"Second test: Successfully published ProductDeleted message",
					logger.Fields{
						"messageId": event.GeMessageId(),
						"productID": event.ProductID,
						"error":     err,
					},
				)

				Convey("It should delete product in the mongo database", func() {
					integrationTestSharedFixture.Log.Infow(
						"Second test: Starting to wait for product deletion",
						logger.Fields{
							"productID": event.ProductID,
						},
					)

					// Wait for the product to be deleted with a more robust condition
					var deletedProduct *models.Product
					var lastErr error
					err := testutils.WaitUntilConditionMet(func() bool {
						deletedProduct, lastErr = integrationTestSharedFixture.ProductRepository.GetProductByProductID(
							ctx,
							integrationTestSharedFixture.Items[0].ProductID,
						)

						integrationTestSharedFixture.Log.Infow(
							"Second test: Checking for product deletion",
							logger.Fields{
								"productID": integrationTestSharedFixture.Items[0].ProductID,
								"deleted":   deletedProduct == nil,
								"error":     lastErr,
							},
						)

						// Return true when product is deleted (nil) or when we get a "not found" error
						return deletedProduct == nil
					}, 30*time.Second) // Increased timeout to 30 seconds

					So(err, ShouldBeNil)
					So(deletedProduct, ShouldBeNil) // Fixed: removed extra parameter
				})
			})
		})

		integrationTestSharedFixture.TearDownTest()
	})

	integrationTestSharedFixture.Log.Info("TearDownSuite started")
	integrationTestSharedFixture.Bus.Stop()
	time.Sleep(1 * time.Second)
}
