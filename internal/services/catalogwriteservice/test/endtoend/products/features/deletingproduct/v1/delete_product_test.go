//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	httpexpect "github.com/gavv/httpexpect/v2"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"
)

var integrationFixture *integration.CatalogWriteIntegrationTestSharedFixture

func TestDeleteProductEndpoint(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewCatalogWriteIntegrationTestSharedFixture(t)
	RunSpecs(t, "DeleteProduct Endpoint EndToEnd Tests")
}

var _ = Describe("DeleteProduct Endpoint", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()

		By("Seeding the required data")
		integrationFixture.SetupTest()
	})

	AfterEach(func() {
		By("Cleanup test data")
		integrationFixture.TearDownTest()
	})

	Describe("Delete product returns appropriate status codes", func() {
		When("An invalid request is made with malformed UUID", func() {
			It("Should return a BadRequest status", func() {
				expect := httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
				expect.DELETE("products/invalid-id").
					WithContext(ctx).
					Expect().
					Status(http.StatusBadRequest)
			})
		})

		When("An invalid request is made with non-existent but valid UUID", func() {
			It("Should return a NotFound status", func() {
				nonExistentID := uuid.New().String()
				expect := httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
				expect.DELETE("products/" + nonExistentID).
					WithContext(ctx).
					Expect().
					Status(http.StatusNotFound)
			})
		})

		When("A valid request is made to delete an existing product", func() {
			It("Should return a NoContent status", func() {
				// First create a product to delete
				expect := httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
				createResponse := expect.POST("products").
					WithContext(ctx).
					WithJSON(map[string]interface{}{
						"name":        "Test Product",
						"description": "Test Description",
						"price":       100.0,
					}).
					Expect().
					Status(http.StatusCreated).
					JSON().
					Object()

				productID := createResponse.Value("productID").String().Raw()

				// Then delete it
				expect.DELETE("products/" + productID).
					WithContext(ctx).
					Expect().
					Status(http.StatusNoContent)
			})
		})
	})
})
