//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	httpexpect "github.com/gavv/httpexpect/v2"
	customTypes "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/customtypes"

	dtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/dtos/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/creatingorder/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/testfixtures/integration"
)

var integrationFixture *integration.OrderIntegrationTestSharedFixture

func TestCreateOrder(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewOrderIntegrationTestSharedFixture(t)
	RunSpecs(t, "CreateOrder Endpoint EndToEnd Tests")
}

var _ = Describe("CreateOrder Feature", func() {
	var (
		ctx     context.Context
		request *dtos.CreateOrderRequestDto
	)

	_ = BeforeEach(func() {
		ctx = context.Background()

		By("Seeding the required data")
		integrationFixture.SetupTest()
	})

	_ = AfterEach(func() {
		By("Cleanup test data")
		integrationFixture.TearDownTest()
	})

	// "Scenario" for testing the creation of an order with valid input
	Describe("Create new order return created status with valid input", func() {
		BeforeEach(func() {
			request = &dtos.CreateOrderRequestDto{
				AccountEmail:    gofakeit.Email(),
				DeliveryAddress: gofakeit.Address().Address,
				DeliveryTime:    customTypes.CustomTime(time.Now()),
				ShopItems: []*dtosV1.ShopItemDto{
					{
						Quantity:    uint64(gofakeit.Number(1, 10)),
						Description: gofakeit.AdjectiveDescriptive(),
						Price:       gofakeit.Price(100, 10000),
						Title:       gofakeit.Name(),
					},
				},
			}
		})
		// "When" step for making a request to create an order
		When("A valid request is made to create an order", func() {
			It("Should returns a StatusCreated response", func() {
				// Create an HTTPExpect instance and make the request
				expect := httpexpect.Default(GinkgoT(), integrationFixture.BaseAddress)
				expect.POST("orders").
					WithContext(ctx).
					WithJSON(request).
					Expect().
					Status(http.StatusCreated)
			})
		})
	})
})
