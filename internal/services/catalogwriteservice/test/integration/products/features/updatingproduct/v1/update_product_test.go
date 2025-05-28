// Package v1 contains the update product test.
package v1

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/gormdbcontext"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/hypothesis"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/messaging"

	mediatr "github.com/mehdihadeli/go-mediatr"
	ginkgo "github.com/onsi/ginkgo"
	gomega "github.com/onsi/gomega"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	uuid "github.com/satori/go.uuid"

	datamodel "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/data/datamodels"
	v1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/updatingproduct/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/updatingproduct/v1/events/integrationevents"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/models"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"
)

var integrationFixture *integration.CatalogWriteIntegrationTestSharedFixture

func TestUpdateProduct(t *testing.T) {
	t.Parallel()
	gomega.RegisterFailHandler(ginkgo.Fail)
	integrationFixture = integration.NewCatalogWriteIntegrationTestSharedFixture(t)
	ginkgo.RunSpecs(t, "Updated Products Integration Tests")
}

var _ = ginkgo.Describe("Update Product Feature", func() {
	// Define variables to hold command and result data
	var (
		ctx             context.Context
		existingProduct *datamodel.ProductDataModel
		command         *v1.UpdateProduct
		result          *mediatr.Unit
		err             error
		id              uuid.UUID
		shouldPublish   hypothesis.Hypothesis[*integrationevents.ProductUpdatedV1]
	)

	_ = ginkgo.BeforeEach(func() {
		ginkgo.By("Seeding the required data")
		integrationFixture.SetupTest()

		existingProduct = integrationFixture.Items[0]
	})

	_ = ginkgo.AfterEach(func() {
		ginkgo.By("Cleanup test data")
		integrationFixture.TearDownTest()
	})

	_ = ginkgo.BeforeSuite(func() {
		ctx = context.Background()

		// in test mode we set rabbitmq `AutoStart=false` in configuration in rabbitmqOptions, so we should run rabbitmq bus manually
		err = integrationFixture.Bus.Start(context.Background())
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

		// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
		time.Sleep(1 * time.Second)
	})

	_ = ginkgo.AfterSuite(func() {
		integrationFixture.Log.Info("TearDownSuite started")
		err := integrationFixture.Bus.Stop()
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		time.Sleep(1 * time.Second)
	})

	// "Scenario" step for testing updating an existing product
	ginkgo.Describe("Updating an existing product in the database", func() {
		ginkgo.Context("Given product exists in the database", func() {
			ginkgo.BeforeEach(func() {
				command, err = v1.NewUpdateProductWithValidation(
					existingProduct.ID,
					"Updated Product ShortTypeName",
					existingProduct.Description,
					existingProduct.Price,
				)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			// "When" step
			ginkgo.When("the UpdateProduct command is executed", func() {
				ginkgo.BeforeEach(func() {
					result, err = mediatr.Send[*v1.UpdateProduct, *mediatr.Unit](
						ctx,
						command,
					)
				})

				// "Then" step
				ginkgo.It("Should not return an error", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(result).NotTo(gomega.BeNil())
				})

				ginkgo.It("Should return a non-nil result", func() {
					gomega.Expect(result).NotTo(gomega.BeNil())
				})

				ginkgo.It(
					"Should update the existing product in the database",
					func() {
						updatedProduct, err := gormdbcontext.FindModelByID[*datamodel.ProductDataModel, *models.Product](
							ctx,
							integrationFixture.CatalogsDBContext,
							existingProduct.ID,
						)
						gomega.Expect(err).To(gomega.BeNil())
						gomega.Expect(updatedProduct).NotTo(gomega.BeNil())
						gomega.Expect(
							updatedProduct.ID,
						).To(gomega.Equal(existingProduct.ID))
						gomega.Expect(
							updatedProduct.Price,
						).To(gomega.Equal(existingProduct.Price))
						gomega.Expect(
							updatedProduct.Name,
						).NotTo(gomega.Equal(existingProduct.Name))
					},
				)
			})
		})
	})

	// "Scenario" step for testing updating a non-existing product
	ginkgo.Describe("Updating a non-existing product in the database", func() {
		ginkgo.Context("Given product not exists in the database", func() {
			ginkgo.BeforeEach(func() {
				// Generate a random ID that does not exist in the database
				id = uuid.NewV4()
				command, err = v1.NewUpdateProductWithValidation(
					id,
					"Updated Product ShortTypeName",
					"Updated Product Description",
					100,
				)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			// "When" step
			ginkgo.When(
				"the UpdateProduct command executed for non-existing product",
				func() {
					ginkgo.BeforeEach(func() {
						result, err = mediatr.Send[*v1.UpdateProduct, *mediatr.Unit](
							ctx,
							command,
						)
					})

					// "Then" step
					ginkgo.It("Should return an error", func() {
						gomega.Expect(err).To(gomega.HaveOccurred())
					})
					ginkgo.It("Should not return a result", func() {
						gomega.Expect(result).To(gomega.BeNil())
					})

					ginkgo.It("Should return a NotFound error", func() {
						gomega.Expect(
							err,
						).To(gomega.MatchError(gomega.ContainSubstring(fmt.Sprintf("product with id `%s` not found", id.String()))))
					})

					ginkgo.It("Should return a custom NotFound error", func() {
						gomega.Expect(customErrors.IsNotFoundError(err)).To(gomega.BeTrue())
						gomega.Expect(
							customErrors.IsApplicationError(
								err,
								http.StatusNotFound,
							),
						).To(gomega.BeTrue())
					})
				},
			)
		})
	})

	// "Scenario" step for testing updating an existing product
	ginkgo.Describe(
		"Publishing ProductUpdated when product updated  successfully",
		func() {
			ginkgo.Context("Given product exists in the database", func() {
				ginkgo.BeforeEach(func() {
					command, err = v1.NewUpdateProductWithValidation(
						existingProduct.ID,
						"Updated Product ShortTypeName",
						existingProduct.Description,
						existingProduct.Price,
					)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					shouldPublish = messaging.ShouldProduced[*integrationevents.ProductUpdatedV1](
						ctx,
						integrationFixture.Bus,
						nil,
					)
				})

				// "When" step
				ginkgo.When(
					"the UpdateProduct command is executed for existing product",
					func() {
						ginkgo.BeforeEach(func() {
							result, err = mediatr.Send[*v1.UpdateProduct, *mediatr.Unit](
								ctx,
								command,
							)
						})

						ginkgo.It(
							"Should publish ProductUpdated event to the broker",
							func() {
								// ensuring message published to the rabbitmq broker
								shouldPublish.Validate(
									ctx,
									"there is no published message",
									time.Second*30,
								)
							},
						)
					},
				)
			})
		},
	)
})
