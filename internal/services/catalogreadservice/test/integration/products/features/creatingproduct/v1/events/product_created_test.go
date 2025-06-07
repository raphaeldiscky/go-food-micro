//go:build integration
// +build integration

package events

// https://github.com/smartystreets/goconvey/wiki

import (
	"context"
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
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
	// in test mode we set rabbitmq `AutoStart=false` in configuration in rabbitmqOptions, so we should run rabbitmq bus manually
	integrationTestSharedFixture.Bus.Start(context.Background())
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
			hypothesis := messaging.ShouldConsume[*externalEvents.ProductCreatedV1](
				ctx,
				integrationTestSharedFixture.Bus,
				nil,
			)

			convey.Convey("When a ProductCreated event consumed", func() {
				err := integrationTestSharedFixture.Bus.PublishMessage(ctx, fakeProduct, nil)
				convey.So(err, convey.ShouldBeNil)

				convey.Convey("Then it should consume the ProductCreated event", func() {
					hypothesis.Validate(ctx, "there is no consumed message", 30*time.Second)
				})
			})
		})

		convey.Convey(
			"Create product in mongo database when a ProductCreated event consumed",
			func() {
				fakeProduct := &externalEvents.ProductCreatedV1{
					Message:     types.NewMessage(uuid.NewV4().String()),
					ProductID:   uuid.NewV4().String(),
					Name:        gofakeit.FirstName(),
					Price:       gofakeit.Price(150, 6000),
					CreatedAt:   time.Now(),
					Description: gofakeit.EmojiDescription(),
				}

				convey.Convey("When a ProductCreated event consumed", func() {
					err := integrationTestSharedFixture.Bus.PublishMessage(ctx, fakeProduct, nil)
					convey.So(err, convey.ShouldBeNil)

					convey.Convey("It should store product in the mongo database", func() {
						ctx := context.Background()
						pid := uuid.NewV4().String()
						productCreated := &externalEvents.ProductCreatedV1{
							Message:     types.NewMessage(uuid.NewV4().String()),
							ProductID:   pid,
							CreatedAt:   time.Now(),
							Name:        gofakeit.Name(),
							Price:       gofakeit.Price(150, 6000),
							Description: gofakeit.AdjectiveDescriptive(),
						}

						err := integrationTestSharedFixture.Bus.PublishMessage(
							ctx,
							productCreated,
							nil,
						)
						convey.So(err, convey.ShouldBeNil)

						var product *models.Product

						err = testutils.WaitUntilConditionMet(func() bool {
							product, err = integrationTestSharedFixture.ProductRepository.GetProductByProductID(
								ctx,
								pid,
							)

							return err == nil && product != nil
						}, 10*time.Second)

						convey.So(err, convey.ShouldBeNil)
						convey.So(product, convey.ShouldNotBeNil)
						convey.So(product.ProductID, convey.ShouldEqual, pid)
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
