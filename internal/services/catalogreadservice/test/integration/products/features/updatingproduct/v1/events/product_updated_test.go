//go:build integration
// +build integration

package events

import (
	"context"
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/messaging"

	. "github.com/smartystreets/goconvey/convey"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	testutils "github.com/raphaeldiscky/go-food-micro/internal/pkg/test/utils"
	uuid "github.com/satori/go.uuid"

	externalEvents "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/updatingproducts/v1/events/integrationevents/externalevents"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/models"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/testfixture/integration"
)

func TestProductUpdatedConsumer(t *testing.T) {
	// Setup and initialization code here.
	integrationTestSharedFixture := integration.NewCatalogReadIntegrationTestSharedFixture(t)

	// in test mode we set rabbitmq `AutoStart=false` in configuration in rabbitmqOptions, so we should run rabbitmq bus manually
	integrationTestSharedFixture.Bus.Start(context.Background())
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)

	Convey("Product Created Feature", t, func() {
		ctx := context.Background()
		integrationTestSharedFixture.SetupTest(t)

		// https://specflow.org/learn/gherkin/#learn-gherkin
		// scenario
		Convey("Consume ProductUpdated event by consumer", func() {
			hypothesis := messaging.ShouldConsume[*externalEvents.ProductUpdatedV1](
				ctx,
				integrationTestSharedFixture.Bus,
				nil,
			)

			fakeUpdateProduct := &externalEvents.ProductUpdatedV1{
				Message:     types.NewMessage(uuid.NewV4().String()),
				ProductID:   integrationTestSharedFixture.Items[0].ProductID,
				Name:        gofakeit.Name(),
				Price:       gofakeit.Price(100, 1000),
				Description: gofakeit.EmojiDescription(),
				UpdatedAt:   time.Now(),
			}

			Convey("When a ProductUpdated event consumed", func() {
				err := integrationTestSharedFixture.Bus.PublishMessage(
					ctx,
					fakeUpdateProduct,
					nil,
				)
				So(err, ShouldBeNil)

				Convey(
					"Then it should consume the ProductUpdated event",
					func() {
						hypothesis.Validate(
							ctx,
							"there is no consumed message",
							30*time.Second,
						)
					},
				)
			})
		})

		// https://specflow.org/learn/gherkin/#learn-gherkin
		// scenario
		Convey(
			"Update product in mongo database when a ProductDeleted event consumed",
			func() {
				fakeUpdateProduct := &externalEvents.ProductUpdatedV1{
					Message:     types.NewMessage(uuid.NewV4().String()),
					ProductID:   integrationTestSharedFixture.Items[0].ProductID,
					Name:        gofakeit.Name(),
					Price:       gofakeit.Price(100, 1000),
					Description: gofakeit.EmojiDescription(),
					UpdatedAt:   time.Now(),
				}

				Convey("When a ProductUpdated event consumed", func() {
					err := integrationTestSharedFixture.Bus.PublishMessage(
						ctx,
						fakeUpdateProduct,
						nil,
					)
					So(err, ShouldBeNil)

					Convey(
						"Then It should update product in the mongo database",
						func() {
							ctx := context.Background()
							productUpdated := &externalEvents.ProductUpdatedV1{
								Message: types.NewMessage(
									uuid.NewV4().String(),
								),
								ProductID:   integrationTestSharedFixture.Items[0].ProductID,
								Name:        gofakeit.Name(),
								Description: gofakeit.AdjectiveDescriptive(),
								Price:       gofakeit.Price(150, 6000),
								UpdatedAt:   time.Now(),
							}

							err := integrationTestSharedFixture.Bus.PublishMessage(
								ctx,
								productUpdated,
								nil,
							)
							So(err, ShouldBeNil)

							var product *models.Product

							err = testutils.WaitUntilConditionMet(func() bool {
								product, err = integrationTestSharedFixture.ProductRepository.GetProductByProductID(
									ctx,
									integrationTestSharedFixture.Items[0].ProductID,
								)

								return product != nil &&
									product.Name == productUpdated.Name
							})

							So(err, ShouldBeNil)
							So(product, ShouldNotBeNil)
							So(
								productUpdated.ProductID,
								ShouldEqual,
								product.ProductID,
							)
						},
					)
				})
			},
		)

		integrationTestSharedFixture.TearDownTest()
	})

	integrationTestSharedFixture.Log.Info("TearDownSuite started")
	integrationTestSharedFixture.Bus.Stop()
	time.Sleep(1 * time.Second)
}
