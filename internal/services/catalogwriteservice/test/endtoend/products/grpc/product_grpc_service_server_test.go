//go:build e2e
// +build e2e

package grpc

import (
	"context"
	"testing"

	"github.com/google/uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	gofakeit "github.com/brianvoe/gofakeit/v6"

	productsservice "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/grpc/genproto"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"
)

var integrationFixture *integration.CatalogWriteIntegrationTestSharedFixture

func TestProductGrpcServiceEndToEnd(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewCatalogWriteIntegrationTestSharedFixture(t)
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
			Name:        gofakeit.Word(),
			Description: gofakeit.AdjectiveDescriptive(),
			Price:       float64(gofakeit.Price(100, 1000)),
		}

		createRes, err := integrationFixture.ProductServiceClient.CreateProduct(ctx, createReq)
		Expect(err).NotTo(HaveOccurred())
		Expect(createRes).NotTo(BeNil())
		Expect(createRes.Id).NotTo(BeEmpty())

		productID = createRes.Id
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
					Name:        gofakeit.Word(),
					Description: gofakeit.AdjectiveDescriptive(),
					Price:       float64(gofakeit.Price(100, 1000)),
				}

				res, err := integrationFixture.ProductServiceClient.CreateProduct(ctx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).NotTo(BeNil())

				By("Verifying the response")
				Expect(res.Id).NotTo(BeEmpty())
				Expect(res.Name).To(Equal(req.Name))
				Expect(res.Description).To(Equal(req.Description))
				Expect(res.Price).To(Equal(req.Price))
				Expect(res.CreatedAt).NotTo(BeEmpty())
				Expect(res.UpdatedAt).NotTo(BeEmpty())
			})
		})

		When("An invalid request is made", func() {
			It("Should return an error for negative price", func() {
				By("Making a request with negative price")
				req := &productsservice.CreateProductReq{
					Name:        gofakeit.Word(),
					Description: gofakeit.AdjectiveDescriptive(),
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
					Description: gofakeit.AdjectiveDescriptive(),
					Price:       float64(gofakeit.Price(100, 1000)),
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
				req := &productsservice.GetProductReq{
					Id: productID,
				}

				res, err := integrationFixture.ProductServiceClient.GetProduct(ctx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).NotTo(BeNil())

				By("Verifying the response")
				Expect(res.Id).To(Equal(productID))
				Expect(res.Name).NotTo(BeEmpty())
				Expect(res.Description).NotTo(BeEmpty())
				Expect(res.Price).To(BeNumerically(">", 0))
				Expect(res.CreatedAt).NotTo(BeEmpty())
				Expect(res.UpdatedAt).NotTo(BeEmpty())
			})
		})

		When("An invalid request is made", func() {
			It("Should return an error for malformed UUID", func() {
				By("Making a request with malformed UUID")
				req := &productsservice.GetProductReq{
					Id: "invalid-uuid",
				}

				res, err := integrationFixture.ProductServiceClient.GetProduct(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(BeNil())
			})

			It("Should return an error for non-existent UUID", func() {
				By("Making a request with non-existent UUID")
				req := &productsservice.GetProductReq{
					Id: uuid.New().String(),
				}

				res, err := integrationFixture.ProductServiceClient.GetProduct(ctx, req)
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
					Id:          productID,
					Name:        gofakeit.Word(),
					Description: gofakeit.AdjectiveDescriptive(),
					Price:       float64(gofakeit.Price(100, 1000)),
				}

				res, err := integrationFixture.ProductServiceClient.UpdateProduct(ctx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).NotTo(BeNil())

				By("Verifying the response")
				Expect(res.Id).To(Equal(productID))
				Expect(res.Name).To(Equal(req.Name))
				Expect(res.Description).To(Equal(req.Description))
				Expect(res.Price).To(Equal(req.Price))
				Expect(res.CreatedAt).NotTo(BeEmpty())
				Expect(res.UpdatedAt).NotTo(BeEmpty())

				By("Verifying the update through a get request")
				getReq := &productsservice.GetProductReq{
					Id: productID,
				}
				getRes, err := integrationFixture.ProductServiceClient.GetProduct(ctx, getReq)
				Expect(err).NotTo(HaveOccurred())
				Expect(getRes).NotTo(BeNil())
				Expect(getRes.Name).To(Equal(req.Name))
				Expect(getRes.Description).To(Equal(req.Description))
				Expect(getRes.Price).To(Equal(req.Price))
			})
		})

		When("An invalid request is made", func() {
			It("Should return an error for malformed UUID", func() {
				By("Making a request with malformed UUID")
				req := &productsservice.UpdateProductReq{
					Id:          "invalid-uuid",
					Name:        gofakeit.Word(),
					Description: gofakeit.AdjectiveDescriptive(),
					Price:       float64(gofakeit.Price(100, 1000)),
				}

				res, err := integrationFixture.ProductServiceClient.UpdateProduct(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(BeNil())
			})

			It("Should return an error for non-existent UUID", func() {
				By("Making a request with non-existent UUID")
				req := &productsservice.UpdateProductReq{
					Id:          uuid.New().String(),
					Name:        gofakeit.Word(),
					Description: gofakeit.AdjectiveDescriptive(),
					Price:       float64(gofakeit.Price(100, 1000)),
				}

				res, err := integrationFixture.ProductServiceClient.UpdateProduct(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(BeNil())
			})

			It("Should return an error for negative price", func() {
				By("Making a request with negative price")
				req := &productsservice.UpdateProductReq{
					Id:          productID,
					Name:        gofakeit.Word(),
					Description: gofakeit.AdjectiveDescriptive(),
					Price:       -100.0,
				}

				res, err := integrationFixture.ProductServiceClient.UpdateProduct(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(BeNil())
			})

			It("Should return an error for empty name", func() {
				By("Making a request with empty name")
				req := &productsservice.UpdateProductReq{
					Id:          productID,
					Name:        "",
					Description: gofakeit.AdjectiveDescriptive(),
					Price:       float64(gofakeit.Price(100, 1000)),
				}

				res, err := integrationFixture.ProductServiceClient.UpdateProduct(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(res).To(BeNil())
			})
		})
	})
})
