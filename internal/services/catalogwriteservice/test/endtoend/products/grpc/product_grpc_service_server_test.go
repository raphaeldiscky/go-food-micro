//go:build e2e
// +build e2e

package grpc

import (
	"context"
	"testing"

	"github.com/google/uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	productsservice "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/grpc/genproto"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"
)

var integrationFixture *integration.CatalogWriteIntegrationTestSharedFixture

func TestProductGrpcServiceEndToEnd(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewCatalogWriteIntegrationTestSharedFixture(t)

	// Wait for gRPC server to be ready
	err := integrationFixture.WaitForGrpcServerReady()
	Expect(err).NotTo(HaveOccurred())

	RunSpecs(t, "Product gRPC Service")
}

var _ = Describe("Product gRPC Service", func() {
	var ctx context.Context
	var productID string

	BeforeEach(func() {
		By("Setting up test context")
		ctx = context.Background()

		By("Seeding the required data")
		integrationFixture.SetupTest()

		// Create a product for testing
		createReq := &productsservice.CreateProductReq{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       10.99,
		}

		createRes, err := integrationFixture.ProductServiceClient.CreateProduct(ctx, createReq)
		Expect(err).NotTo(HaveOccurred())
		Expect(createRes).NotTo(BeNil())
		Expect(createRes.ProductID).NotTo(BeEmpty())

		productID = createRes.ProductID
	})

	AfterEach(func() {
		By("Cleanup test data")
		integrationFixture.TearDownTest()
	})

	Describe("Create product", func() {
		When("A valid request is made", func() {
			It("Should create a product successfully", func() {
				By("Making a request to create a product")
				req := &productsservice.CreateProductReq{
					Name:        "New Test Product",
					Description: "New Test Description",
					Price:       20.99,
				}

				res, err := integrationFixture.ProductServiceClient.CreateProduct(ctx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).NotTo(BeNil())

				By("Verifying the response")
				Expect(res.ProductID).NotTo(BeEmpty())

				// Verify through GetProductByID
				getReq := &productsservice.GetProductByIDReq{
					ProductID: res.ProductID,
				}
				getRes, err := integrationFixture.ProductServiceClient.GetProductByID(ctx, getReq)
				Expect(err).NotTo(HaveOccurred())
				Expect(getRes).NotTo(BeNil())
				Expect(getRes.Product).NotTo(BeNil())
				Expect(getRes.Product.Name).To(Equal(req.Name))
				Expect(getRes.Product.Description).To(Equal(req.Description))
				Expect(getRes.Product.Price).To(Equal(req.Price))
			})
		})

		When("An invalid request is made", func() {
			It("Should return an error for negative price", func() {
				By("Making a request with negative price")
				req := &productsservice.CreateProductReq{
					Name:        "Test Product",
					Description: "Test Description",
					Price:       -100.0,
				}

				res, err := integrationFixture.ProductServiceClient.CreateProduct(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(BeNil())
			})

			It("Should return an error for empty name", func() {
				By("Making a request with empty name")
				req := &productsservice.CreateProductReq{
					Name:        "",
					Description: "Test Description",
					Price:       10.99,
				}

				res, err := integrationFixture.ProductServiceClient.CreateProduct(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(BeNil())
			})
		})
	})

	Describe("Get product", func() {
		When("A valid request is made", func() {
			It("Should return the product successfully", func() {
				By("Making a request to get the product")
				req := &productsservice.GetProductByIDReq{
					ProductID: productID,
				}

				res, err := integrationFixture.ProductServiceClient.GetProductByID(ctx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).NotTo(BeNil())
				Expect(res.Product).NotTo(BeNil())

				By("Verifying the response")
				Expect(res.Product.ProductID).To(Equal(productID))
				Expect(res.Product.Name).NotTo(BeEmpty())
				Expect(res.Product.Description).NotTo(BeEmpty())
				Expect(res.Product.Price).To(BeNumerically(">", 0))
				Expect(res.Product.CreatedAt).NotTo(BeNil())
				Expect(res.Product.UpdatedAt).NotTo(BeNil())
			})
		})

		When("An invalid request is made", func() {
			It("Should return an error for malformed UUID", func() {
				By("Making a request with malformed UUID")
				req := &productsservice.GetProductByIDReq{
					ProductID: "invalid-uuid",
				}

				res, err := integrationFixture.ProductServiceClient.GetProductByID(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(BeNil())
			})

			It("Should return an error for non-existent UUID", func() {
				By("Making a request with non-existent UUID")
				req := &productsservice.GetProductByIDReq{
					ProductID: uuid.New().String(),
				}

				res, err := integrationFixture.ProductServiceClient.GetProductByID(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(BeNil())
			})
		})
	})

	Describe("Update product", func() {
		When("A valid request is made", func() {
			It("Should update the product successfully", func() {
				By("Making a request to update the product")
				req := &productsservice.UpdateProductReq{
					ProductID:   productID,
					Name:        "Updated Test Product",
					Description: "Updated Test Description",
					Price:       30.99,
				}

				res, err := integrationFixture.ProductServiceClient.UpdateProduct(ctx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).NotTo(BeNil())

				// Since UpdateProductRes is empty, we verify the update through a get request
				By("Verifying the update through a get request")
				getReq := &productsservice.GetProductByIDReq{
					ProductID: productID,
				}
				getRes, err := integrationFixture.ProductServiceClient.GetProductByID(ctx, getReq)
				Expect(err).NotTo(HaveOccurred())
				Expect(getRes).NotTo(BeNil())
				Expect(getRes.Product).NotTo(BeNil())
				Expect(getRes.Product.ProductID).To(Equal(productID))
				Expect(getRes.Product.Name).To(Equal(req.Name))
				Expect(getRes.Product.Description).To(Equal(req.Description))
				Expect(getRes.Product.Price).To(Equal(req.Price))
			})
		})

		When("An invalid request is made", func() {
			It("Should return an error for malformed UUID", func() {
				By("Making a request with malformed UUID")
				req := &productsservice.UpdateProductReq{
					ProductID:   "invalid-uuid",
					Name:        "Updated Test Product",
					Description: "Updated Test Description",
					Price:       30.99,
				}

				res, err := integrationFixture.ProductServiceClient.UpdateProduct(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(BeNil())
			})

			It("Should return an error for non-existent UUID", func() {
				By("Making a request with non-existent UUID")
				req := &productsservice.UpdateProductReq{
					ProductID:   uuid.New().String(),
					Name:        "Updated Test Product",
					Description: "Updated Test Description",
					Price:       30.99,
				}

				res, err := integrationFixture.ProductServiceClient.UpdateProduct(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(BeNil())
			})

			It("Should return an error for negative price", func() {
				By("Making a request with negative price")
				req := &productsservice.UpdateProductReq{
					ProductID:   productID,
					Name:        "Updated Test Product",
					Description: "Updated Test Description",
					Price:       -100.0,
				}

				res, err := integrationFixture.ProductServiceClient.UpdateProduct(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(BeNil())
			})

			It("Should return an error for empty name", func() {
				By("Making a request with empty name")
				req := &productsservice.UpdateProductReq{
					ProductID:   productID,
					Name:        "",
					Description: "Updated Test Description",
					Price:       30.99,
				}

				res, err := integrationFixture.ProductServiceClient.UpdateProduct(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(BeNil())
			})
		})
	})
})
