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
	testUtils "github.com/raphaeldiscky/go-food-micro/internal/pkg/test/utils"
	externalEvents "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/updating_products/v1/events/integration_events/external_events"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/models"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/testfixture/integration"

	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProductUpdatedConsumer(t *testing.T) {
	// Setup and initialization code here.
	integrationTestSharedFixture := integration.NewIntegrationTestSharedFixture(
		t,
	)

	// Start the bus and wait for it to be fully ready
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	integrationTestSharedFixture.Log.Info("Starting RabbitMQ bus...")
	err := integrationTestSharedFixture.Bus.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start bus: %v", err)
	}
	integrationTestSharedFixture.Log.Info("RabbitMQ bus started successfully")

	// Wait longer for consumers to be ready
	integrationTestSharedFixture.Log.Info("Waiting for consumers to be ready...")
	time.Sleep(10 * time.Second)
	integrationTestSharedFixture.Log.Info("Consumers should be ready now")

	Convey("Product Updated Feature", t, func() {
		ctx := context.Background()
		integrationTestSharedFixture.Log.Info("Setting up test data...")
		integrationTestSharedFixture.SetupTest(t)

		// Verify test data setup
		if len(integrationTestSharedFixture.Items) == 0 {
			t.Fatal("No test items were created during setup")
		}
		testProduct := integrationTestSharedFixture.Items[0]
		if testProduct.ProductId == "" {
			t.Fatal("Test product ID is empty")
		}
		integrationTestSharedFixture.Log.Infow(
			"Test data setup complete",
			logger.Fields{
				"productId": testProduct.ProductId,
				"name":      testProduct.Name,
				"price":     testProduct.Price,
			},
		)

		// https://specflow.org/learn/gherkin/#learn-gherkin
		// scenario
		Convey("Consume ProductUpdated event by consumer", func() {
			// Create a test product update event
			fakeUpdateProduct := &externalEvents.ProductUpdatedV1{
				Message:     types.NewMessage(uuid.NewV4().String()),
				ProductId:   testProduct.ProductId, // Use the actual product ID from test data
				Name:        gofakeit.Name(),
				Price:       gofakeit.Price(100, 1000),
				Description: gofakeit.EmojiDescription(),
				UpdatedAt:   time.Now(),
			}

			integrationTestSharedFixture.Log.Infow(
				"Created test product update event",
				logger.Fields{
					"productId":   fakeUpdateProduct.ProductId,
					"name":        fakeUpdateProduct.Name,
					"price":       fakeUpdateProduct.Price,
					"description": fakeUpdateProduct.Description,
				},
			)

			// Create a channel to signal when the message is consumed
			messageConsumed := make(chan bool, 1)
			messageError := make(chan error, 1)

			// Create hypothesis and validate it
			integrationTestSharedFixture.Log.Info("Setting up message consumer hypothesis...")
			_ = messaging.ShouldConsume[*externalEvents.ProductUpdatedV1](
				ctx,
				integrationTestSharedFixture.Bus,
				func(msg *externalEvents.ProductUpdatedV1) bool {
					if msg == nil {
						err := fmt.Errorf("received nil message")
						integrationTestSharedFixture.Log.Errorw("Message consumption error", logger.Fields{"error": err})
						messageError <- err
						return false
					}

					integrationTestSharedFixture.Log.Infow(
						"Received message",
						logger.Fields{
							"productId":   msg.ProductId,
							"name":        msg.Name,
							"price":       msg.Price,
							"description": msg.Description,
						},
					)

					if msg.ProductId == fakeUpdateProduct.ProductId &&
						msg.Name == fakeUpdateProduct.Name &&
						msg.Price == fakeUpdateProduct.Price &&
						msg.Description == fakeUpdateProduct.Description {
						integrationTestSharedFixture.Log.Info("Message matches expected values")
						messageConsumed <- true
						return true
					}

					integrationTestSharedFixture.Log.Infow(
						"Message does not match expected values",
						logger.Fields{
							"expected": fakeUpdateProduct,
							"actual":   msg,
						},
					)
					return false
				},
			)
			integrationTestSharedFixture.Log.Info("Message consumer hypothesis setup complete")

			Convey("When a ProductUpdated event is published", func() {
				// Publish the message with retry
				integrationTestSharedFixture.Log.Info("Attempting to publish message...")
				var publishErr error
				for i := 0; i < 3; i++ {
					integrationTestSharedFixture.Log.Infow(
						"Publishing message attempt",
						logger.Fields{"attempt": i + 1},
					)
					publishErr = integrationTestSharedFixture.Bus.PublishMessage(
						ctx,
						fakeUpdateProduct,
						nil,
					)
					if publishErr == nil {
						integrationTestSharedFixture.Log.Info("Message published successfully")
						break
					}
					integrationTestSharedFixture.Log.Warn(
						"Failed to publish message, retrying...",
						logger.Fields{"error": publishErr},
					)
					time.Sleep(time.Second)
				}
				So(publishErr, ShouldBeNil)

				Convey("Then it should consume the ProductUpdated event and update the database", func() {
					integrationTestSharedFixture.Log.Info("Waiting for message consumption...")
					// Wait for message consumption with timeout
					select {
					case err := <-messageError:
						integrationTestSharedFixture.Log.Errorw(
							"Message consumption failed",
							logger.Fields{"error": err},
						)
						t.Fatalf("Error consuming message: %v", err)
					case <-messageConsumed:
						integrationTestSharedFixture.Log.Info("Message consumed successfully, waiting for database update...")
						// Message was consumed, now wait for database update
						var product *models.Product
						var dbErr error
						err = testUtils.WaitUntilConditionMet(func() bool {
							product, dbErr = integrationTestSharedFixture.ProductRepository.GetProductByProductId(
								ctx,
								fakeUpdateProduct.ProductId,
							)
							if dbErr != nil {
								integrationTestSharedFixture.Log.Errorw(
									"Error getting product from database",
									logger.Fields{"error": dbErr},
								)
								return false
							}

							if product == nil {
								integrationTestSharedFixture.Log.Info("Product not found in database yet")
								return false
							}

							matches := product.Name == fakeUpdateProduct.Name &&
								product.Description == fakeUpdateProduct.Description &&
								product.Price == fakeUpdateProduct.Price

							if !matches {
								integrationTestSharedFixture.Log.Infow(
									"Product in database doesn't match expected values",
									logger.Fields{
										"expected": fakeUpdateProduct,
										"actual":   product,
									},
								)
							} else {
								integrationTestSharedFixture.Log.Info("Product in database matches expected values")
							}

							return matches
						}, 45*time.Second) // Increased timeout for database update

						if err != nil {
							integrationTestSharedFixture.Log.Errorw(
								"Database update timeout",
								logger.Fields{"error": err},
							)
						}
						if dbErr != nil {
							integrationTestSharedFixture.Log.Errorw(
								"Database error",
								logger.Fields{"error": dbErr},
							)
						}

						So(err, ShouldBeNil, "Database update timeout")
						So(dbErr, ShouldBeNil, "Database error")
						So(product, ShouldNotBeNil, "Product not found in database")
						So(product.ProductId, ShouldEqual, fakeUpdateProduct.ProductId)
						So(product.Name, ShouldEqual, fakeUpdateProduct.Name)
						So(product.Description, ShouldEqual, fakeUpdateProduct.Description)
						So(product.Price, ShouldEqual, fakeUpdateProduct.Price)
					case <-time.After(45 * time.Second): // Increased timeout for message consumption
						integrationTestSharedFixture.Log.Error("Message consumption timeout - the event was not consumed within the expected time")
						t.Fatal("Message consumption timeout - the event was not consumed within the expected time")
					}
				})
			})
		})

		integrationTestSharedFixture.Log.Info("Starting test teardown...")
		integrationTestSharedFixture.TearDownTest()
		integrationTestSharedFixture.Log.Info("Test teardown complete")
	})

	integrationTestSharedFixture.Log.Info("Starting suite teardown...")
	integrationTestSharedFixture.Bus.Stop()
	time.Sleep(5 * time.Second) // Increased wait time after stopping bus
	integrationTestSharedFixture.Log.Info("Suite teardown complete")
}
