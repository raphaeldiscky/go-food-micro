//go:build integration
// +build integration

package events

import (
	"context"
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/messaging"

	. "github.com/smartystreets/goconvey/convey"

	testUtils "github.com/raphaeldiscky/go-food-micro/internal/pkg/test/utils"
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
	time.Sleep(1 * time.Second)

	Convey("Product Deleted Feature", t, func() {
		ctx := context.Background()
		// will execute with each subtest
		integrationTestSharedFixture.SetupTest(t)

		Convey("Consume ProductDeleted event by consumer", func() {
			event := &externalevents.ProductDeletedV1{
				Message:   types.NewMessage(uuid.NewV4().String()),
				ProductID: integrationTestSharedFixture.Items[0].ProductID,
			}
			// check for consuming `ProductDeletedV1` message with existing consumer
			hypothesis := messaging.ShouldConsume[*externalevents.ProductDeletedV1](
				ctx,
				integrationTestSharedFixture.Bus,
				nil,
			)

			Convey("When a ProductDeleted event consumed", func() {
				err := integrationTestSharedFixture.Bus.PublishMessage(
					ctx,
					event,
					nil,
				)
				So(err, ShouldBeNil)

				Convey("Then it should consume the ProductDeleted event", func() {
					hypothesis.Validate(ctx, "there is no consumed message", 30*time.Second)
				})
			})
		})

		// https://specflow.org/learn/gherkin/#learn-gherkin
		// scenario
		Convey("Delete product in mongo database when a ProductDeleted event consumed", func() {
			event := &externalevents.ProductDeletedV1{
				Message:   types.NewMessage(uuid.NewV4().String()),
				ProductID: integrationTestSharedFixture.Items[0].ProductID,
			}

			// First verify the product exists
			existingProduct, err := integrationTestSharedFixture.ProductRepository.GetProductByProductID(
				ctx,
				integrationTestSharedFixture.Items[0].ProductID,
			)
			So(err, ShouldBeNil)
			So(existingProduct, ShouldNotBeNil)

			Convey("When a ProductDeleted event consumed", func() {
				err := integrationTestSharedFixture.Bus.PublishMessage(
					ctx,
					event,
					nil,
				)
				So(err, ShouldBeNil)

				Convey("It should delete product in the mongo database", func() {
					// Wait for the product to be deleted with a more robust condition
					var deletedProduct *models.Product
					err := testUtils.WaitUntilConditionMet(func() bool {
						deletedProduct, err = integrationTestSharedFixture.ProductRepository.GetProductByProductID(
							ctx,
							integrationTestSharedFixture.Items[0].ProductID,
						)
						if err != nil {
							integrationTestSharedFixture.Log.Errorw(
								"Error checking product deletion",
								logger.Fields{
									"error":     err,
									"productID": integrationTestSharedFixture.Items[0].ProductID,
								},
							)
							return false
						}
						return deletedProduct == nil
					}, 10*time.Second)

					So(err, ShouldBeNil)
					So(deletedProduct, ShouldBeNil, "Product should be deleted")
				})
			})
		})

		integrationTestSharedFixture.TearDownTest()
	})

	integrationTestSharedFixture.Log.Info("TearDownSuite started")
	integrationTestSharedFixture.Bus.Stop()
	time.Sleep(1 * time.Second)
}
