//go:build integration
// +build integration

package events

// https://github.com/smartystreets/goconvey/wiki

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/messaging"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	testutils "github.com/raphaeldiscky/go-food-micro/internal/pkg/test/utils"
	uuid "github.com/satori/go.uuid"
	convey "github.com/smartystreets/goconvey/convey"

	externalEvents "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/creatingproduct/v1/events/integrationevents/externalevents"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/models"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/testfixture/integration"
)

// TestProductCreatedConsumer is a test for the ProductCreated consumer
func TestProductCreatedConsumer(t *testing.T) {
	// Setup and initialization code here.
	integrationTestSharedFixture := integration.NewCatalogReadIntegrationTestSharedFixture(t)
	// The bus is already started by the integration fixture, no need to start it again
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)

	convey.Convey("Product Created Feature", t, func() {
		// will execute with each subtest
		integrationTestSharedFixture.SetupTest(t)
		ctx := context.Background()

		// https://specflow.org/learn/gherkin/#learn-gherkin
		// scenario
		convey.Convey("Consume ProductCreated event by consumer", func() {
			fakeProduct := &externalEvents.ProductCreatedV1{
				Message:     types.NewMessage(uuid.NewV4().String()),
				ProductID:   uuid.NewV4().String(),
				Name:        gofakeit.FirstName(),
				Price:       gofakeit.Price(150, 6000),
				CreatedAt:   time.Now(),
				Description: gofakeit.EmojiDescription(),
			}

			// Add extensive debug logging
			integrationTestSharedFixture.Log.Infow(
				"Creating hypothesis for message consumption",
				logger.Fields{
					"messageType":     fmt.Sprintf("%T", fakeProduct),
					"messageTypeName": fakeProduct.GetMessageTypeName(),
				},
			)

			hypothesis := messaging.ShouldConsume[*externalEvents.ProductCreatedV1](
				ctx,
				integrationTestSharedFixture.Bus,
				nil,
			)

			convey.Convey("When a ProductCreated event consumed", func() {
				integrationTestSharedFixture.Log.Infow(
					"About to publish ProductCreated message",
					logger.Fields{
						"messageId":       fakeProduct.GeMessageId(),
						"messageType":     fmt.Sprintf("%T", fakeProduct),
						"messageTypeName": fakeProduct.GetMessageTypeName(),
						"productID":       fakeProduct.ProductID,
						"name":            fakeProduct.Name,
						"price":           fakeProduct.Price,
					},
				)

				err := integrationTestSharedFixture.Bus.PublishMessage(ctx, fakeProduct, nil)
				convey.So(err, convey.ShouldBeNil)

				integrationTestSharedFixture.Log.Infow(
					"Successfully published ProductCreated message",
					logger.Fields{
						"messageId": fakeProduct.GeMessageId(),
						"error":     err,
					},
				)

				convey.Convey("Then it should consume the ProductCreated event", func() {
					integrationTestSharedFixture.Log.Infow(
						"Starting hypothesis validation with 30 second timeout",
						logger.Fields{},
					)
					hypothesis.Validate(ctx, "there is no consumed message", 30*time.Second)
				})
			})
		})

		convey.Convey(
			"Create product in mongo database when a ProductCreated event consumed",
			func() {
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

				// Create a single ProductCreated event
				pid := uuid.NewV4().String()
				productCreated := &externalEvents.ProductCreatedV1{
					Message:     types.NewMessage(uuid.NewV4().String()),
					ProductID:   pid,
					CreatedAt:   time.Now(),
					Name:        gofakeit.Name(),
					Price:       gofakeit.Price(150, 6000),
					Description: gofakeit.AdjectiveDescriptive(),
				}

				integrationTestSharedFixture.Log.Infow(
					"Second test: About to publish ProductCreated message",
					logger.Fields{
						"messageId":       productCreated.GeMessageId(),
						"messageType":     fmt.Sprintf("%T", productCreated),
						"messageTypeName": productCreated.GetMessageTypeName(),
						"productID":       productCreated.ProductID,
						"name":            productCreated.Name,
						"price":           productCreated.Price,
					},
				)

				convey.Convey("When a ProductCreated event consumed", func() {
					// Publish the message
					err := integrationTestSharedFixture.Bus.PublishMessage(
						ctx,
						productCreated,
						nil,
					)
					convey.So(err, convey.ShouldBeNil)

					integrationTestSharedFixture.Log.Infow(
						"Second test: Successfully published ProductCreated message",
						logger.Fields{
							"messageId": productCreated.GeMessageId(),
							"productID": productCreated.ProductID,
							"error":     err,
						},
					)

					convey.Convey("It should store product in the mongo database", func() {
						var product *models.Product

						integrationTestSharedFixture.Log.Infow(
							"Second test: Starting to wait for product in database",
							logger.Fields{
								"productID": pid,
							},
						)

						// Wait for the message to be consumed and the product to be created in the database
						err = testutils.WaitUntilConditionMet(func() bool {
							product, err = integrationTestSharedFixture.ProductRepository.GetProductByProductID(
								ctx,
								pid,
							)

							integrationTestSharedFixture.Log.Infow(
								"Second test: Checking for product in database",
								logger.Fields{
									"productID": pid,
									"found":     product != nil,
									"error":     err,
								},
							)

							return err == nil && product != nil
						}, 15*time.Second)

						convey.So(err, convey.ShouldBeNil)
						convey.So(product, convey.ShouldNotBeNil)
						convey.So(product.ProductID, convey.ShouldEqual, pid)
						convey.So(product.Name, convey.ShouldEqual, productCreated.Name)
						convey.So(
							product.Description,
							convey.ShouldEqual,
							productCreated.Description,
						)
						convey.So(product.Price, convey.ShouldEqual, productCreated.Price)
					})
				})
			},
		)

		integrationTestSharedFixture.TearDownTest()
	})

	integrationTestSharedFixture.Log.Info("TearDownSuite started")
	integrationTestSharedFixture.Bus.Stop()
	time.Sleep(1 * time.Second)
}
