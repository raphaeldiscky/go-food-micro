//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	httpexpect "github.com/gavv/httpexpect/v2"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/testfixtures/integration"
)

var integrationFixture *integration.OrderIntegrationTestSharedFixture

func TestGetProductById(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewOrderIntegrationTestSharedFixture(t)
	RunSpecs(t, "GetOrderByID Endpoint EndToEnd Tests")
}

var _ = Describe("GetOrderByID Feature", func() {
	var (
		ctx context.Context
		id  string
	)

	_ = BeforeEach(func() {
		ctx = context.Background()

		By("Seeding the required data")
		integrationFixture.SetupTest()

		id = integrationFixture.Items[0].ID
	})

	_ = AfterEach(func() {
		By("Cleanup test data")
		integrationFixture.TearDownTest()
	})

	// "Scenario" for testing the retrieval of an order by a valid ID
	Describe("Get order by ID with a valid ID returns ok status", func() {
		// "When" step for making a request to get an order by ID
		When("A valid request is made with a valid ID", func() {
			It("should return an 'OK' status", func() {
				expect := httpexpect.Default(GinkgoT(), integrationFixture.BaseAddress)
				expect.GET("orders/{id}").
					WithPath("id", id).
					WithContext(ctx).
					Expect().
					Status(http.StatusOK)
			})
		})
	})
})
