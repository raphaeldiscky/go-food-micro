//go:build integration
// +build integration

package v1

import (
	"context"
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/hypothesis"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/messaging"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	mediatr "github.com/mehdihadeli/go-mediatr"
	testUtils "github.com/raphaeldiscky/go-food-micro/internal/pkg/test/utils"

	dtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/dtos/v1"
	createOrderCommandV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/creatingorder/v1/commands"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/creatingorder/v1/dtos"
	integrationEvents "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/creatingorder/v1/events/integrationevents"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/readmodels"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/testfixtures/integration"
)

var integrationFixture *integration.OrderIntegrationTestSharedFixture

func TestCreateOrder(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewOrderIntegrationTestSharedFixture(t)
	RunSpecs(t, "Create Order Integration Tests")
}

var _ = Describe("Create Order Feature", func() {
	var (
		ctx          context.Context
		err          error
		command      *createOrderCommandV1.CreateOrder
		result       *dtos.CreateOrderResponseDto
		createdOrder *readmodels.OrderReadModel
		// id            string
		shouldPublish hypothesis.Hypothesis[*integrationEvents.OrderCreatedV1]
	)

	_ = BeforeEach(func() {
		By("Seeding the required data")
		integrationFixture.SetupTest()

		// id = integrationFixture.Items[0].OrderID
	})

	_ = AfterEach(func() {
		By("Cleanup test data")
		integrationFixture.TearDownTest()
	})

	_ = BeforeSuite(func() {
		ctx = context.Background()

		// in test mode we set rabbitmq `AutoStart=false` in configuration in rabbitmqOptions, so we should run rabbitmq bus manually
		err = integrationFixture.Bus.Start(context.Background())
		Expect(err).ShouldNot(HaveOccurred())

		// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
		time.Sleep(1 * time.Second)
	})

	_ = AfterSuite(func() {
		integrationFixture.Log.Info("TearDownSuite started")
		err := integrationFixture.Bus.Stop()
		Expect(err).ShouldNot(HaveOccurred())
		time.Sleep(1 * time.Second)
	})

	// "Scenario" for testing the creation of a new order
	Describe("Creating a new order in EventStoreDB", func() {
		BeforeEach(func() {
			command, err = createOrderCommandV1.NewCreateOrder(
				[]*dtosV1.ShopItemDto{
					{
						Quantity:    uint64(gofakeit.Number(1, 10)),
						Description: gofakeit.AdjectiveDescriptive(),
						Price:       gofakeit.Price(100, 10000),
						Title:       gofakeit.Name(),
					},
				},
				gofakeit.Email(),
				gofakeit.Address().Address,
				time.Now(),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(command).ToNot(BeNil())
		})
		When("the CreateOrder command is executed for non-existing order", func() {
			BeforeEach(func() {
				// "When" step for executing the createOrderCommand
				result, err = mediatr.Send[*createOrderCommandV1.CreateOrder, *dtos.CreateOrderResponseDto](
					ctx,
					command,
				)
			})
			// "Then" step for expected behavior
			It("Should create the order successfully", func() {
				// "Then" step for assertions
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(command.OrderID).To(Equal(result.OrderID))
			})
		})
	})

	// "Scenario" for testing the creation of a new order in MongoDB Read
	Describe("Creating a new order in MongoDB Read", func() {
		BeforeEach(func() {
			command, err = createOrderCommandV1.NewCreateOrder(
				[]*dtosV1.ShopItemDto{
					{
						Quantity:    uint64(gofakeit.Number(1, 10)),
						Description: gofakeit.AdjectiveDescriptive(),
						Price:       gofakeit.Price(100, 10000),
						Title:       gofakeit.Name(),
					},
				},
				gofakeit.Email(),
				gofakeit.Address().Address,
				time.Now(),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(command).ToNot(BeNil())
		})
		// "When" step for creating a new order
		When("the CreateOrder command is executed for non-existing order", func() {
			BeforeEach(func() {
				// "When" step for executing the createOrderCommand
				result, err = mediatr.Send[*createOrderCommandV1.CreateOrder, *dtos.CreateOrderResponseDto](
					context.Background(),
					command,
				)
			})

			It("Should create the order successfully", func() {
				// "Then" step for assertions
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
			})

			// "Then" step for expected behavior
			It("Should retrieve created order in MongoDB Read database", func() {
				// Use a utility function to wait until the order is available in MongoDB Read
				err = testUtils.WaitUntilConditionMet(func() bool {
					createdOrder, err = integrationFixture.OrderMongoRepository.GetOrderByOrderID(
						ctx,
						result.OrderID,
					)
					Expect(err).ToNot(HaveOccurred())
					return createdOrder != nil
				})

				Expect(err).To(BeNil())
			})
		})
	})

	// "Scenario" for testing the publishing of an "OrderCreated" event
	Describe(
		"Publishing an OrderCreated event to the message broker when order saved successfully",
		func() {
			BeforeEach(func() {
				command, err = createOrderCommandV1.NewCreateOrder(
					[]*dtosV1.ShopItemDto{
						{
							Quantity:    uint64(gofakeit.Number(1, 10)),
							Description: gofakeit.AdjectiveDescriptive(),
							Price:       gofakeit.Price(100, 10000),
							Title:       gofakeit.Name(),
						},
					},
					gofakeit.Email(),
					gofakeit.Address().Address,
					time.Now(),
				)

				Expect(err).ToNot(HaveOccurred())
				Expect(command).ToNot(BeNil())
			})

			// "When" step for creating and sending an order
			When("CreateOrder command is executed for non-existing order", func() {
				BeforeEach(func() {
					// Ensure shouldPublish is properly initialized before executing the command
					shouldPublish = messaging.ShouldProduced[*integrationEvents.OrderCreatedV1](
						ctx,
						integrationFixture.Bus, nil,
					)

					// "When" step for executing the createOrderCommand
					result, err = mediatr.Send[*createOrderCommandV1.CreateOrder, *dtos.CreateOrderResponseDto](
						context.Background(),
						command,
					)
				})

				It("Should return no error", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("Should return not nil result", func() {
					Expect(result).ToNot(BeNil())
				})

				It("Should publish OrderCreated event to the broker", func() {
					// ensuring message published to the rabbitmq broker
					shouldPublish.Validate(ctx, "there is no published message", time.Second*30)
				})
			})
		},
	)
})
