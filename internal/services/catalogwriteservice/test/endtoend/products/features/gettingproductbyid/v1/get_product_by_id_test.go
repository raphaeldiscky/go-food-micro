//go:build e2e
// +build e2e

package v1

import (
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gofakeit "github.com/brianvoe/gofakeit/v6"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"
)

var integrationFixture *integration.CatalogWriteIntegrationTestSharedFixture

func TestGetProductByIdEndToEnd(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewCatalogWriteIntegrationTestSharedFixture(t)
	RunSpecs(t, "GetProductById Endpoint")
}

var _ = Describe("GetProductById Endpoint", func() {
	var expect *httpexpect.Expect
	var productID string

	BeforeEach(func() {
		By("Seeding the required data")
		integrationFixture.SetupTest()

		expect = httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)

		// Create a product for testing
		createRequest := map[string]interface{}{
			"name":        gofakeit.Name(),
			"description": gofakeit.AdjectiveDescriptive(),
			"price":       float64(gofakeit.Price(100, 1000)),
		}

		// Create the product and verify the response
		createRes := expect.POST("/products").
			WithJSON(createRequest).
			Expect().
			Status(201).
			JSON().
			Object()

		// Ensure we got a valid productID
		createRes.ContainsKey("productID")
		productID = createRes.Value("productID").String().Raw()
		Expect(productID).NotTo(BeEmpty(), "Product ID should not be empty")

		// Verify it's a valid UUID by attempting to parse it
		_, err := uuid.Parse(productID)
		Expect(err).NotTo(HaveOccurred(), "Product ID should be a valid UUID")

		// Verify the product exists by getting it
		res := expect.GET("/products/{id}", productID).
			Expect().
			Status(200).
			JSON().
			Object()

		// Verify the response structure with nested product object
		res.ContainsKey("product")
		product := res.Value("product").Object()
		product.ContainsKey("id")
		product.Value("id").String().Equal(productID)
	})

	AfterEach(func() {
		By("Cleanup test data")
		integrationFixture.TearDownTest()
	})

	Describe("Get product by ID", func() {
		When("A valid product ID is provided", func() {
			It("Should return the product details", func() {
				By("Getting the product by ID")
				res := expect.GET("/products/{id}", productID).
					Expect().
					Status(200).
					JSON().
					Object()

				By("Verifying the response structure")
				res.ContainsKey("product")
				product := res.Value("product").Object()
				product.ContainsKey("id")
				product.ContainsKey("name")
				product.ContainsKey("description")
				product.ContainsKey("price")
				product.ContainsKey("createdAt")
				product.ContainsKey("updatedAt")

				By("Verifying the product ID matches")
				product.Value("id").String().Equal(productID)
			})
		})

		When("An invalid UUID is provided", func() {
			It("Should return a 400 Bad Request", func() {
				expect.GET("/products/{id}", "invalid-uuid").
					Expect().
					Status(400)
			})
		})

		When("A non-existent product ID is provided", func() {
			It("Should return a 404 Not Found", func() {
				nonExistentID := uuid.New().String()
				expect.GET("/products/{id}", nonExistentID).
					Expect().
					Status(404)
			})
		})
	})
})
