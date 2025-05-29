//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	httpexpect "github.com/gavv/httpexpect/v2"
	ginkgo "github.com/onsi/ginkgo"
	gomega "github.com/onsi/gomega"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"
)

var integrationFixture *integration.IntegrationTestSharedFixture

func TestCreateProductEndpoint(t *testing.T) {
	ginkgo.RegisterFailHandler(gomega.Fail)
	integrationFixture = integration.NewIntegrationTestSharedFixture(t)
	ginkgo.RunSpecs(t, "CreateProduct Endpoint EndToEnd Tests")
}

var _ = ginkgo.Describe("CreateProduct Feature", func() {
	var (
		ctx     context.Context
		request *dtos.CreateProductRequestDto
	)

	_ = ginkgo.BeforeEach(func() {
		ctx = context.Background()

		By("Seeding the required data")
		integrationFixture.SetupTest()
	})

	_ = ginkgo.AfterEach(func() {
		By("Cleanup test data")
		integrationFixture.TearDownTest()
	})

	// "Scenario" step for testing the create product API with valid input
	ginkgo.Describe("Create new product return created status with valid input", func() {
		ginkgo.BeforeEach(func() {
			// Generate a valid request
			request = &dtos.CreateProductRequestDto{
				Description: gofakeit.AdjectiveDescriptive(),
				Price:       gofakeit.Price(100, 1000),
				Name:        gofakeit.Name(),
			}
		})
		// "When" step
		ginkgo.When("A valid request is made to create a product", func() {
			// "Then" step
			ginkgo.It("Should returns a StatusCreated response", func() {
				// Create an HTTPExpect instance and make the request
				expect := httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
				expect.POST("products").
					WithContext(ctx).
					WithJSON(request).
					Expect().
					Status(http.StatusCreated)
			})
		})
	})

	// "Scenario" step for testing the create product API with invalid price input
	ginkgo.Describe("Create product returns a BadRequest status with invalid price input", func() {
		ginkgo.BeforeEach(func() {
			// Generate an invalid request with zero price
			request = &dtos.CreateProductRequestDto{
				Description: gofakeit.AdjectiveDescriptive(),
				Price:       0.0,
				Name:        gofakeit.Name(),
			}
		})
		// "When" step
		ginkgo.When("An invalid request is made with a zero price", func() {
			// "Then" step
			ginkgo.It("Should return a BadRequest status", func() {
				// Create an HTTPExpect instance and make the request
				expect := httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
				expect.POST("products").
					WithContext(ctx).
					WithJSON(request).
					Expect().
					Status(http.StatusBadRequest)
			})
		})
	})
})
