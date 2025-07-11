//go:build integration
// +build integration

package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/models"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"
)

var integrationFixture *integration.CatalogWriteIntegrationTestSharedFixture

func TestProductPostgresRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewCatalogWriteIntegrationTestSharedFixture(t)
	RunSpecs(t, "ProductPostgresRepository Integration Tests")
}

var _ = Describe("Product Repository Suite", func() {
	// Define variables to hold repository and product data
	var (
		ctx             context.Context
		product         *models.Product
		createdProduct  *models.Product
		updatedProduct  *models.Product
		existingProduct *models.Product
		err             error
		id              uuid.UUID
	)

	_ = BeforeEach(func() {
		By("Seeding the required data")
		integrationFixture.SetupTest()

		id = integrationFixture.Items[0].ID
	})

	_ = AfterEach(func() {
		By("Cleanup test data")
		integrationFixture.TearDownTest()
	})

	_ = BeforeSuite(func() {
		ctx = context.Background()

		// in test mode we set rabbitmq `AutoStart=false` in configuration in rabbitmqOptions, so we should run rabbitmq bus manually
		err = integrationFixture.Bus.Start(context.Background())
		Expect(err).ShouldNot(HaveOccurred())

		// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
		time.Sleep(1 * time.Second)
	})

	_ = AfterSuite(func() {
		integrationFixture.Log.Info("TearDownSuite started")
		err := integrationFixture.Bus.Stop()
		Expect(err).ShouldNot(HaveOccurred())
		time.Sleep(1 * time.Second)
	})

	// "Scenario" step for testing creating a new product in the database
	Describe("Creating a new product in the database", func() {
		BeforeEach(func() {
			product = &models.Product{
				Name:        gofakeit.Name(),
				Description: gofakeit.AdjectiveDescriptive(),
				ID:          uuid.NewV4(),
				Price:       gofakeit.Price(100, 1000),
				CreatedAt:   time.Now(),
			}
		})

		// "When" step
		When("CreateProduct function of ProductRepository executed", func() {
			BeforeEach(func() {
				createdProduct, err = integrationFixture.ProductRepository.CreateProduct(
					ctx,
					product,
				)
			})

			// "Then" step
			It("Should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("Should return a non-nil created product", func() {
				Expect(createdProduct).NotTo(BeNil())
			})

			It("Should have the same ID as the input product", func() {
				Expect(createdProduct.ID).To(Equal(product.ID))
			})

			It("Should be able to retrieve the created product from the database", func() {
				retrievedProduct, err := integrationFixture.ProductRepository.GetProductByID(
					ctx,
					createdProduct.ID,
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(retrievedProduct).NotTo(BeNil())
				Expect(retrievedProduct.ID).To(Equal(createdProduct.ID))
			})
		})
	})

	// "Scenario" step for testing updating an existing product in the database
	Describe("Updating an existing product in the database", func() {
		BeforeEach(func() {
			existingProduct, err = integrationFixture.ProductRepository.GetProductByID(ctx, id)
			Expect(err).To(BeNil())
			Expect(existingProduct).NotTo(BeNil())
		})

		// "When" step
		When("UpdateProduct function of ProductRepository executed", func() {
			BeforeEach(func() {
				// Update the name of the existing product
				existingProduct.Name = "Updated Product ShortTypeName"
				_, err = integrationFixture.ProductRepository.UpdateProduct(ctx, existingProduct)
			})

			// "Then" step
			It("Should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("Should be able to retrieve the updated product from the database", func() {
				updatedProduct, err = integrationFixture.ProductRepository.GetProductByID(
					ctx,
					existingProduct.ID,
				)
				Expect(err).To(BeNil())
				Expect(updatedProduct).NotTo(BeNil())
				Expect(updatedProduct.Name).To(Equal("Updated Product ShortTypeName"))
				// You can add more assertions to validate other properties of the updated product
			})
		})
	})

	// "Scenario" step for testing deleting an existing product in the database
	Describe("Deleting an existing product from the database", func() {
		BeforeEach(func() {
			// Ensure that the product with 'id' exists in the database
			product, err := integrationFixture.ProductRepository.GetProductByID(ctx, id)
			Expect(err).To(BeNil())
			Expect(product).NotTo(BeNil())
		})

		// "When" step
		When("DeleteProduct function of ProductRepository executed", func() {
			BeforeEach(func() {
				err = integrationFixture.ProductRepository.DeleteProductByID(ctx, id)
			})

			// "Then" step
			It("Should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("Should delete given product from the database", func() {
				product, err := integrationFixture.ProductRepository.GetProductByID(ctx, id)
				Expect(err).To(HaveOccurred())
				Expect(customErrors.IsNotFoundError(err)).To(BeTrue())
				Expect(product).To(BeNil())
			})
		})
	})

	// "Scenario" step for testing retrieving an existing product from the database
	Describe("Retrieving an existing product from the database", func() {
		BeforeEach(func() {
			// Ensure that the product with 'id' exists in the database
			product, err := integrationFixture.ProductRepository.GetProductByID(ctx, id)
			Expect(err).To(BeNil())
			Expect(product).NotTo(BeNil())
		})

		// "When" step
		When("GetProductByID function of ProductRepository executed", func() {
			BeforeEach(func() {
				existingProduct, err = integrationFixture.ProductRepository.GetProductByID(ctx, id)
			})
			It("should not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(existingProduct).NotTo(BeNil())
			})

			It("should retrieve correct data from database by ID", func() {
				Expect(existingProduct.ID).To(Equal(id))
			})
		})
	})

	// "Scenario" step for testing retrieving a product that does not exist in the database
	Describe("Retrieving a product that does not exist in the database", func() {
		BeforeEach(func() {
			// Ensure that the product with 'id' exists in the database
			product, err := integrationFixture.ProductRepository.GetProductByID(ctx, id)
			Expect(err).To(BeNil())
			Expect(product).NotTo(BeNil())
		})

		// "When" step
		When("GetProductByID function of ProductRepository executed", func() {
			BeforeEach(func() {
				// Use a random UUID that does not exist in the database
				nonexistentID := uuid.NewV4()
				existingProduct, err = integrationFixture.ProductRepository.GetProductByID(
					ctx,
					nonexistentID,
				)
			})

			// "Then" step
			It("Should return a NotFound error", func() {
				Expect(err).To(HaveOccurred())
				Expect(customErrors.IsNotFoundError(err)).To(BeTrue())
			})

			It("Should not return a product", func() {
				Expect(existingProduct).To(BeNil())
			})
		})
	})

	// "Scenario" step for testing retrieving all existing products from the database
	Describe("Retrieving all existing products from the database", func() {
		// "When" step
		When("GetAllProducts function of ProductRepository executed", func() {
			It("should not return an error and return the correct number of products", func() {
				res, err := integrationFixture.ProductRepository.GetAllProducts(
					ctx,
					utils.NewListQuery(10, 1),
				)
				Expect(err).To(BeNil())
				Expect(res).NotTo(BeNil())
				Expect(len(res.Items)).To(Equal(2)) // Replace with the expected number of products
			})
		})
	})
})
