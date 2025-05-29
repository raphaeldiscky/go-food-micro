//go:build integration
// +build integration

package v1

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mediatr "github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/gettingorderbyid/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/gettingorderbyid/v1/queries"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/testfixtures/integration"
)

var integrationFixture *integration.IntegrationTestSharedFixture

func TestGetOrderByID(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewIntegrationTestSharedFixture(t)
	RunSpecs(t, "Get Order By ID Integration Tests")
}

var _ = Describe("Get Order By ID Feature", func() {
	var (
		ctx    context.Context
		query  *queries.GetOrderByID
		err    error
		id     uuid.UUID
		result *dtos.GetOrderByIDResponseDto
	)

	_ = BeforeEach(func() {
		By("Seeding the required data")
		integrationFixture.SetupTest()

		idString := integrationFixture.Items[0].ID
		id, err = uuid.FromString(idString)
		Expect(err).NotTo(HaveOccurred())
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

	// "Scenario" for testing the retrieval of an existing order by ID
	Describe("Retrieving an existing order by ID from the database", func() {
		BeforeEach(func() {
			query, err = queries.NewGetOrderByID(id)
			Expect(err).ToNot(HaveOccurred())
			Expect(query).ToNot(BeNil())
		})
		// "When" step for executing the query to get the order
		When("retrieving an existing order by ID", func() {
			BeforeEach(func() {
				// "When" step for executing the query to get the order
				result, err = mediatr.Send[*queries.GetOrderByID, *dtos.GetOrderByIDResponseDto](
					ctx,
					query,
				)
			})

			It("Should return the order successfully with correct properties", func() {
				// "Then" step for assertions
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Order).NotTo(BeNil())
				Expect(result.Order.ID).To(Equal(id.String()))
				Expect(result.Order.OrderID).NotTo(BeEmpty())
			})
		})
	})
})
