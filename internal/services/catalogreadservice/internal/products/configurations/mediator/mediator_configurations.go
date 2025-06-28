// Package mediator contains the mediator configurations.
package mediator

import (
	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"

	mediatr "github.com/mehdihadeli/go-mediatr"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/contracts/data"
	v1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/creatingproduct/v1"
	createProductDtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/creatingproduct/v1/dtos"
	deleteProductCommandV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/deletingproducts/v1/commands"
	getProductByIdDtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/getproductbyid/v1/dtos"
	getProductByIdQueryV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/getproductbyid/v1/queries"
	getProductsDtoV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/gettingproducts/v1/dtos"
	getProductsQueryV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/gettingproducts/v1/queries"
	searchProductsDtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/searchingproducts/v1/dtos"
	searchProductsQueryV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/searchingproducts/v1/queries"
	updateProductCommandV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/updatingproducts/v1/commands"
)

// ConfigProductsMediator configures the products mediator.
func ConfigProductsMediator(
	log logger.Logger,
	mongoProductRepository data.ProductRepository,
	cacheProductRepository data.ProductCacheRepository,
	tracer tracing.AppTracer,
) error {
	log.Infow("Starting mediator configuration for products", logger.Fields{})

	err := mediatr.RegisterRequestHandler[*v1.CreateProduct, *createProductDtosV1.CreateProductResponseDto](
		v1.NewCreateProductHandler(
			log,
			mongoProductRepository,
			cacheProductRepository,
			tracer,
		),
	)
	if err != nil {
		log.Errorw("Failed to register CreateProduct handler", logger.Fields{"error": err})

		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}
	log.Infow("Successfully registered CreateProduct handler", logger.Fields{})

	err = mediatr.RegisterRequestHandler[*deleteProductCommandV1.DeleteProduct, *mediatr.Unit](
		deleteProductCommandV1.NewDeleteProductHandler(
			log,
			mongoProductRepository,
			cacheProductRepository,
			tracer,
		),
	)
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*updateProductCommandV1.UpdateProduct, *mediatr.Unit](
		updateProductCommandV1.NewUpdateProductHandler(
			log,
			mongoProductRepository,
			cacheProductRepository,
			tracer,
		),
	)
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*getProductsQueryV1.GetProducts, *getProductsDtoV1.GetProductsResponseDto](
		getProductsQueryV1.NewGetProductsHandler(log, mongoProductRepository, tracer),
	)
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*searchProductsQueryV1.SearchProducts, *searchProductsDtosV1.SearchProductsResponseDto](
		searchProductsQueryV1.NewSearchProductsHandler(
			log,
			mongoProductRepository,
			tracer,
		),
	)
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	err = mediatr.RegisterRequestHandler[*getProductByIdQueryV1.GetProductByID, *getProductByIdDtosV1.GetProductByIDResponseDto](
		getProductByIdQueryV1.NewGetProductByIDHandler(
			log,
			mongoProductRepository,
			cacheProductRepository,
			tracer,
		),
	)
	if err != nil {
		return errors.WrapIf(err, "error while registering handlers in the mediator")
	}

	log.Infow(
		"Successfully completed mediator configuration for all product handlers",
		logger.Fields{},
	)

	return nil
}
