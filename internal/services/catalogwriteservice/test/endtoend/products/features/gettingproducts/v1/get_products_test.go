//go:build e2e
// +build e2e

package v1

import (
	"testing"

	"github.com/gavv/httpexpect/v2"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gofakeit "github.com/brianvoe/gofakeit/v6"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"
)

var integrationFixture *integration.CatalogWriteIntegrationTestSharedFixture

func TestGetProductsEndToEnd(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewCatalogWriteIntegrationTestSharedFixture(t)
	RunSpecs(t, "GetProducts Endpoint")
}

var _ = Describe("GetProducts Endpoint", func() {
	var expect *httpexpect.Expect

	BeforeEach(func() {
		By("Seeding the required data")
		integrationFixture.SetupTest()

		expect = httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
	})

	AfterEach(func() {
		By("Cleanup test data")
		integrationFixture.TearDownTest()
	})

	Describe("Get products returns OK status", func() {
		When("A request is made to get all products", func() {
			It("Should return an OK status with products data", func() {
				// Create a product first to ensure we have data
				createRequest := map[string]interface{}{
					"name":        gofakeit.Name(),
					"description": gofakeit.AdjectiveDescriptive(),
					"price":       float64(gofakeit.Price(100, 1000)),
				}

				createRes := expect.POST("/products").
					WithJSON(createRequest).
					Expect().
					Status(201).
					JSON().
					Object()

				createRes.ContainsKey("productID")

				// Get all products
				res := expect.GET("/products").
					Expect().
					Status(200).
					JSON().
					Object()

				// Verify the response structure
				res.ContainsKey("Products")
				products := res.Value("Products").Object()
				products.ContainsKey("items")
				products.ContainsKey("page")
				products.ContainsKey("size")

				// Verify the items array
				items := products.Value("items").Array()
				items.Length().Gt(0)

				// Verify the first item structure
				firstItem := items.First().Object()
				firstItem.ContainsKey("id")
				firstItem.ContainsKey("name")
				firstItem.ContainsKey("description")
				firstItem.ContainsKey("price")
				firstItem.ContainsKey("createdAt")
				firstItem.ContainsKey("updatedAt")
			})
		})
	})
})
