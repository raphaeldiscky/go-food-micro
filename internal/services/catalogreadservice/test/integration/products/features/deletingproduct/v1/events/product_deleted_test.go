//go:build integration
// +build integration

package events

import (
	"context"
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"

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
	// in test mode we set rabbitmq `AutoStart=false` in configuration in rabbitmqOptions, so we should run rabbitmq bus manually
	integrationTestSharedFixture.Bus.Start(context.Background())
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(2 * time.Second) // Increased wait time for consumer readiness

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

		Convey("Delete product in mongo database when a ProductDeleted event consumed", func() {
			event := &externalevents.ProductDeletedV1{
				Message:   types.NewMessage(uuid.NewV4().String()),
				ProductID: integrationTestSharedFixture.Items[0].ProductID,
			}

			Convey("When a ProductDeleted event consumed", func() {
				integrationTestSharedFixture.Log.Infow(
					"Publishing ProductDeleted event",
					logger.Fields{
						"event":     event,
						"productID": event.ProductID,
						"messageID": event.Message,
					},
				)

				err := integrationTestSharedFixture.Bus.PublishMessage(
					ctx,
					event,
					nil,
				)
				So(err, ShouldBeNil)

				Convey("It should delete product in the mongo database", func() {
					// Wait for the product to be deleted with a more robust condition
					var deletedProduct *models.Product
					var lastErr error
					err := testutils.WaitUntilConditionMet(func() bool {
						deletedProduct, lastErr = integrationTestSharedFixture.ProductRepository.GetProductByProductID(
							ctx,
							integrationTestSharedFixture.Items[0].ProductID,
						)
						if lastErr != nil {
							integrationTestSharedFixture.Log.Errorw(
								"Error checking product deletion",
								logger.Fields{
									"error":     lastErr,
									"productID": integrationTestSharedFixture.Items[0].ProductID,
								},
							)
							return false
						}
						if deletedProduct != nil {
							integrationTestSharedFixture.Log.Infow(
								"Product still exists, waiting for deletion",
								logger.Fields{
									"productID": deletedProduct.ProductID,
									"name":      deletedProduct.Name,
								},
							)
						}
						return deletedProduct == nil
					}, 30*time.Second) // Increased timeout to 30 seconds
					if err != nil {
						integrationTestSharedFixture.Log.Errorw(
							"Timeout waiting for product deletion",
							logger.Fields{
								"error":     err,
								"lastError": lastErr,
								"productID": integrationTestSharedFixture.Items[0].ProductID,
							},
						)
					}

					So(err, ShouldBeNil)
					So(deletedProduct, ShouldBeNil, "Product should be deleted")
				})
			})
		})

		integrationTestSharedFixture.TearDownTest()
	})

	integrationTestSharedFixture.Log.Info("TearDownSuite started")
	integrationTestSharedFixture.Bus.Stop()
	time.Sleep(2 * time.Second) // Increased wait time for cleanup
}
