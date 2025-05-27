//go:build integration
// +build integration

package events

import (
	"context"
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	testUtils "github.com/raphaeldiscky/go-food-micro/internal/pkg/test/utils"
	externalEvents "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/deleting_products/v1/events/integration_events/external_events"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/models"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/testfixture/integration"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProductDeleted(t *testing.T) {
	// Setup and initialization code here.
	integrationTestSharedFixture := integration.NewIntegrationTestSharedFixture(t)
	require.NotNil(t, integrationTestSharedFixture, "Integration test shared fixture should not be nil")

	Convey("Product Deleted Feature", t, func() {
		ctx := context.Background()
		// will execute with each subtest
		integrationTestSharedFixture.SetupTest(t)

		// https://specflow.org/learn/gherkin/#learn-gherkin
		// scenario
		Convey("Delete product in mongo database when a ProductDeleted event consumed", func() {
			event := &externalEvents.ProductDeletedV1{
				Message:   types.NewMessage(uuid.NewV4().String()),
				ProductId: integrationTestSharedFixture.Items[0].ProductId,
			}

			// First verify the product exists
			existingProduct, err := integrationTestSharedFixture.ProductRepository.GetProductByProductId(
				ctx,
				integrationTestSharedFixture.Items[0].ProductId,
			)
			So(err, ShouldBeNil)
			So(existingProduct, ShouldNotBeNil)

			Convey("When a ProductDeleted event consumed", func() {
				// Publish the message with retries
				var publishErr error
				for i := 0; i < 3; i++ {
					publishErr = integrationTestSharedFixture.Bus.PublishMessage(
						ctx,
						event,
						nil,
					)
					if publishErr == nil {
						break
					}
					integrationTestSharedFixture.Log.Warn(
						"Failed to publish message, retrying...",
						logger.Fields{
							"attempt": i + 1,
							"error":   publishErr,
						},
					)
					time.Sleep(time.Second)
				}
				So(publishErr, ShouldBeNil)

				Convey("It should delete product in the mongo database", func() {
					// Wait for the product to be deleted with a more robust condition
					var deletedProduct *models.Product
					err := testUtils.WaitUntilConditionMet(func() bool {
						deletedProduct, err = integrationTestSharedFixture.ProductRepository.GetProductByProductId(
							ctx,
							integrationTestSharedFixture.Items[0].ProductId,
						)
						if err != nil {
							integrationTestSharedFixture.Log.Errorw(
								"Error checking product deletion",
								logger.Fields{
									"error":     err,
									"productId": integrationTestSharedFixture.Items[0].ProductId,
								},
							)
							return false
						}
						return deletedProduct == nil
					}, 45*time.Second) // Increased timeout to 45 seconds

					So(err, ShouldBeNil, "Timeout waiting for product deletion")
					So(deletedProduct, ShouldBeNil, "Product should be deleted")
				})
			})
		})

		integrationTestSharedFixture.TearDownTest()
	})

	integrationTestSharedFixture.Log.Info("TearDownSuite started")
}
