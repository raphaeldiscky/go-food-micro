//go:build integration
// +build integration

package data

import (
	"context"
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	. "github.com/smartystreets/goconvey/convey"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/models"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/testfixture/integration"
)

func TestProductPostgresRepository(t *testing.T) {
	integrationTestSharedFixture := integration.NewCatalogReadIntegrationTestSharedFixture(t)

	// scenario
	Convey("MongoDB Product Repository", t, func() {
		ctx := context.Background()
		integrationTestSharedFixture.SetupTest(t)

		Convey("When we create the new product in the database", func() {
			product := &models.Product{
				ID:          uuid.NewV4().String(),
				ProductID:   uuid.NewV4().String(),
				Name:        gofakeit.Name(),
				Description: gofakeit.AdjectiveDescriptive(),
				Price:       gofakeit.Price(100, 1000),
				CreatedAt:   time.Now(),
			}

			createdProduct, err := integrationTestSharedFixture.ProductRepository.CreateProduct(
				ctx,
				product,
			)

			Convey("Then the product should be created successfully", func() {
				// Assert that there is no error during creation.
				So(err, ShouldBeNil)

				Convey("And we should be able to retrieve the product by ID", func() {
					retrievedProduct, err := integrationTestSharedFixture.ProductRepository.GetProductByID(
						ctx,
						createdProduct.ID,
					)

					Convey("And retrieved product should match the created product", func() {
						// Assert that there is no error during retrieval.
						So(err, ShouldBeNil)

						// Assert that the retrieved product matches the created product.
						So(retrievedProduct.ID, ShouldEqual, createdProduct.ID)
					})
				})
			})
		})

		Convey("When we delete the existing product", func() {
			id := integrationTestSharedFixture.Items[0].ID
			err := integrationTestSharedFixture.ProductRepository.DeleteProductByID(ctx, id)

			Convey("Then the product should be deleted successfully", func() {
				// Ensure there is no error during deletion.
				So(err, ShouldBeNil)

				Convey("And when we attempt to retrieve the product by ID", func() {
					product, err := integrationTestSharedFixture.ProductRepository.GetProductByID(
						ctx,
						id,
					)

					Convey("And error should occur indicating the product is not found", func() {
						// Verify that there is an error.
						So(err, ShouldNotBeNil)

						// Check if the error is of a specific type (e.g., a not found error).
						So(customErrors.IsNotFoundError(err), ShouldBeTrue)

						// Verify that the retrieved product is nil.
						So(product, ShouldBeNil)
					})
				})
			})
		})

		Convey("When we update the existing product", func() {
			Convey("Then the product should be updated successfully", func() {
				id := integrationTestSharedFixture.Items[0].ID
				existingProduct, err := integrationTestSharedFixture.ProductRepository.GetProductByID(
					ctx,
					id,
				)

				// Make sure the existing product exists and there is no error.
				So(err, ShouldBeNil)
				So(existingProduct, ShouldNotBeNil)

				// Modify the existing product's name.
				existingProduct.Name = "test_update_product"

				// Update the product in the database.
				_, err = integrationTestSharedFixture.ProductRepository.UpdateProduct(
					ctx,
					existingProduct,
				)

				// Ensure there is no error during the update.
				So(err, ShouldBeNil)

				// Retrieve the updated product from the database.
				updatedProduct, err := integrationTestSharedFixture.ProductRepository.GetProductByID(
					ctx,
					id,
				)
				So(err, ShouldBeNil)

				// Verify that the updated product's name matches the new name.
				So(updatedProduct.Name, ShouldEqual, "test_update_product")
			})
		})

		Convey("When attempting to get a product that does not exist", func() {
			res, err := integrationTestSharedFixture.ProductRepository.GetProductByID(
				ctx,
				uuid.NewV4().String(),
			)

			Convey("Then it should return a NotFound error and nil result", func() {
				// Verify that there is an error.
				So(err, ShouldNotBeNil)

				// Check if the error is of a specific type (e.g., a not found error).
				So(customErrors.IsNotFoundError(err), ShouldBeTrue)

				// Verify that the retrieved result is nil.
				So(res, ShouldBeNil)
			})
		})

		Convey("When attempting to get an existing product from the database", func() {
			id := integrationTestSharedFixture.Items[0].ID
			res, err := integrationTestSharedFixture.ProductRepository.GetProductByID(ctx, id)

			Convey("Then it should return the product and no error", func() {
				// Ensure there is no error.
				So(err, ShouldBeNil)

				// Verify that the result is not nil.
				So(res, ShouldNotBeNil)

				// Verify that the retrieved product's ID matches the expected ID.
				So(res.ID, ShouldEqual, id)
			})
		})

		Convey("When attempting to get all existing products from the database", func() {
			res, err := integrationTestSharedFixture.ProductRepository.GetAllProducts(
				ctx,
				utils.NewListQuery(10, 1),
			)

			Convey("Then it should return the list of products and no error", func() {
				// Ensure there is no error.
				So(err, ShouldBeNil)

				// Verify the expected number of products in the list (2 products are seeded in SetupTest).
				So(len(res.Items), ShouldEqual, 2)
			})
		})

		integrationTestSharedFixture.TearDownTest()
	})
}
