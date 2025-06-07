//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	httpexpect "github.com/gavv/httpexpect/v2"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"
)

var integrationFixture *integration.CatalogWriteIntegrationTestSharedFixture

func TestCreateProductEndpoint(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewCatalogWriteIntegrationTestSharedFixture(t)
	RunSpecs(t, "CreateProduct Endpoint EndToEnd Tests")
}

var _ = Describe("CreateProduct Endpoint", func() {
	var (
		ctx     context.Context
		request *dtos.CreateProductRequestDto
	)

	BeforeEach(func() {
		ctx = context.Background()

		By("Seeding the required data")
		integrationFixture.SetupTest()
	})

	AfterEach(func() {
		By("Cleanup test data")
		integrationFixture.TearDownTest()
	})

	Describe("Create new product return created status with valid input", func() {
		BeforeEach(func() {
			// Generate a valid request with explicit float64 price
			price := float64(gofakeit.Price(100, 1000))
			request = &dtos.CreateProductRequestDto{
				Description: gofakeit.AdjectiveDescriptive(),
				Price:       price,
				Name:        gofakeit.Name(),
			}
		})

		When("A valid request is made to create a product", func() {
			It("Should returns a StatusCreated response", func() {
				// Create an HTTPExpect instance and make the request
				expect := httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
				obj := expect.POST("products").
					WithContext(ctx).
					WithJSON(map[string]interface{}{
						"name":        request.Name,
						"description": request.Description,
						"price":       request.Price,
					}).
					Expect().
					Status(http.StatusCreated).
					JSON().
					Object()

				// Verify response structure
				obj.ContainsKey("productID")
				Expect(obj.Value("productID").Raw()).NotTo(BeEmpty())
			})
		})
	})

	Describe("Create product returns a BadRequest status with invalid price input", func() {
		BeforeEach(func() {
			// Generate an invalid request with zero price
			request = &dtos.CreateProductRequestDto{
				Description: gofakeit.AdjectiveDescriptive(),
				Price:       0.0,
				Name:        gofakeit.Name(),
			}
		})

		When("An invalid request is made with a zero price", func() {
			It("Should return a BadRequest status", func() {
				// Create an HTTPExpect instance and make the request
				expect := httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
				expect.POST("products").
					WithContext(ctx).
					WithJSON(map[string]interface{}{
						"name":        request.Name,
						"description": request.Description,
						"price":       0.0,
					}).
					Expect().
					Status(http.StatusBadRequest)
			})
		})
	})
})
