//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"
)

var integrationFixture *integration.CatalogWriteIntegrationTestSharedFixture

func TestUpdateProductEndToEnd(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewCatalogWriteIntegrationTestSharedFixture(t)
	RunSpecs(t, "UpdateProduct Endpoint")
}

var _ = Describe("UpdateProduct Endpoint", func() {
	var (
		ctx       context.Context
		expect    *httpexpect.Expect
		productID uuid.UUID
	)

	BeforeEach(func() {
		By("Setting up test data")
		ctx = context.Background()
		integrationFixture.SetupTest()
		expect = httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
		productID = integrationFixture.Items[0].ID
	})

	AfterEach(func() {
		By("Cleaning up test data")
		integrationFixture.TearDownTest()
	})

	Describe("Update product returns appropriate status codes", func() {
		When("A valid update request is made", func() {
			It("Should return a 204 No Content status", func() {
				By("Making update request")
				updateRequest := map[string]interface{}{
					"name":        gofakeit.Name(),
					"description": gofakeit.AdjectiveDescriptive(),
					"price":       gofakeit.Price(100, 1000),
				}

				// First verify the product exists
				expect.GET("/products/{id}", productID).
					WithContext(ctx).
					Expect().
					Status(http.StatusOK)

				// Then attempt the update
				expect.PUT("/products/{id}", productID).
					WithContext(ctx).
					WithJSON(updateRequest).
					Expect().
					Status(http.StatusNoContent)

				// Verify the update through a get request
				response := expect.GET("/products/{id}", productID).
					WithContext(ctx).
					Expect().
					Status(http.StatusOK).
					JSON().
					Object()

				By("Verifying updated values")
				product := response.Value("product").Object()
				product.Value("id").String().Equal(productID.String())
				product.Value("name").String().Equal(updateRequest["name"].(string))
				product.Value("description").String().Equal(updateRequest["description"].(string))
				product.Value("price").Number().Equal(updateRequest["price"].(float64))
			})

			It("Should return a 400 Bad Request for invalid UUID", func() {
				By("Making request with invalid UUID")
				invalidUUID := "not-a-uuid"
				updateRequest := map[string]interface{}{
					"name":        gofakeit.Name(),
					"description": gofakeit.AdjectiveDescriptive(),
					"price":       gofakeit.Price(100, 1000),
				}

				expect.PUT("/products/{id}", invalidUUID).
					WithContext(ctx).
					WithJSON(updateRequest).
					Expect().
					Status(http.StatusBadRequest)
			})

			It("Should return a 404 Not Found for non-existent product", func() {
				By("Making request with non-existent UUID")
				nonExistentID := uuid.NewV4()
				updateRequest := map[string]interface{}{
					"name":        gofakeit.Name(),
					"description": gofakeit.AdjectiveDescriptive(),
					"price":       gofakeit.Price(100, 1000),
				}

				expect.PUT("/products/{id}", nonExistentID).
					WithContext(ctx).
					WithJSON(updateRequest).
					Expect().
					Status(http.StatusNotFound)
			})

			It("Should return a 400 Bad Request for invalid data", func() {
				By("Making request with invalid data")
				invalidRequest := map[string]interface{}{
					"name":        "", // Empty name should fail validation
					"description": gofakeit.AdjectiveDescriptive(),
					"price":       0, // Zero price should fail validation
				}

				expect.PUT("/products/{id}", productID).
					WithContext(ctx).
					WithJSON(invalidRequest).
					Expect().
					Status(http.StatusBadRequest)
			})
		})
	})
})
