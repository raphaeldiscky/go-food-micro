//go:build e2e
// +build e2e

package grpc

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"

	productService "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/grpc/genproto"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"
)

var integrationFixture *integration.CatalogWriteIntegrationTestSharedFixture

func TestProductGrpcServiceEndToEnd(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewIntegrationTestSharedFixture(t)
	RunSpecs(t, "ProductGrpcService EndToEnd Tests")
}

var _ = Describe("Product Grpc Service Feature", func() {
	var (
		ctx context.Context
		id  uuid.UUID
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

	// "Scenario" step for testing the creation of a product with valid data in the database
	Describe("Creation of a product with valid data in the database", func() {
		// "When" step
		When("A request is made to create a product with valid data", func() {
			// "Then" step
			It("Should return a non-empty ID", func() {
				// Create a gRPC request with valid data
				request := &productService.CreateProductReq{
					Price:       gofakeit.Price(100, 1000),
					Name:        gofakeit.Name(),
					Description: gofakeit.AdjectiveDescriptive(),
				}

				// Make the gRPC request to create the product
				res, err := integrationFixture.ProductServiceClient.CreateProduct(ctx, request)
				Expect(err).To(BeNil())
				Expect(res).NotTo(BeNil())
				Expect(res.ProductID).NotTo(BeEmpty())
			})
		})
	})

	// "Scenario" step for testing the retrieval of data with a valid ID
	Describe("Retrieve product with a valid ID", func() {
		// "When" step
		When("A request is made to retrieve data with a valid ID", func() {
			// "Then" step
			It("Should return data with a matching ID", func() {
				// Make the gRPC request to retrieve data by ID
				res, err := integrationFixture.ProductServiceClient.GetProductByID(
					ctx,
					&productService.GetProductByIDReq{ProductID: id.String()},
				)

				Expect(err).To(BeNil())
				Expect(res).NotTo(BeNil())
				Expect(res.Product).NotTo(BeNil())
				Expect(res.Product.ProductID).To(Equal(id.String()))
			})
		})
	})
})
