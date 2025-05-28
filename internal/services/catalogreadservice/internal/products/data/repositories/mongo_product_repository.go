package repositories

// https://github.com/Kamva/mgm
// https://github.com/mongodb/mongo-go-driver
// https://blog.logrocket.com/how-to-use-mongodb-with-go/
// https://www.mongodb.com/docs/drivers/go/current/quick-reference/
// https://www.mongodb.com/docs/drivers/go/current/fundamentals/bson/
// https://www.mongodb.com/docs

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/data"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb/repository"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/attribute"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"

	utils2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/utils"
	uuid2 "github.com/satori/go.uuid"
	attribute2 "go.opentelemetry.io/otel/attribute"

	data2 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/contracts/data"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/models"
)

const (
	productCollection = "products"
)

// mongoProductRepository is a struct that contains the mongo product repository.
type mongoProductRepository struct {
	log                    logger.Logger
	mongoGenericRepository data.GenericRepository[*models.Product]
	tracer                 tracing.AppTracer
}

// NewMongoProductRepository creates a new MongoProductRepository.
func NewMongoProductRepository(
	log logger.Logger,
	db *mongo.Client,
	mongoOptions *mongodb.MongoDbOptions,
	tracer tracing.AppTracer,
) data2.ProductRepository {
	mongoRepo := repository.NewGenericMongoRepository[*models.Product](
		db,
		mongoOptions.Database,
		productCollection,
	)

	return &mongoProductRepository{
		log:                    log,
		mongoGenericRepository: mongoRepo,
		tracer:                 tracer,
	}
}

// GetAllProducts gets all products from the database.
func (p *mongoProductRepository) GetAllProducts(
	ctx context.Context,
	listQuery *utils.ListQuery,
) (*utils.ListResult[*models.Product], error) {
	ctx, span := p.tracer.Start(ctx, "mongoProductRepository.GetAllProducts")
	defer span.End()

	// https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/query-document/
	result, err := p.mongoGenericRepository.GetAll(ctx, listQuery)
	if err != nil {
		return nil, utils2.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"error in the paginate",
			),
		)
	}

	p.log.Infow(
		"products loaded",
		logger.Fields{"ProductsResult": result},
	)

	span.SetAttributes(attribute.Object("ProductsResult", result))

	return result, nil
}

// SearchProducts searches for products in the database.
func (p *mongoProductRepository) SearchProducts(
	ctx context.Context,
	searchText string,
	listQuery *utils.ListQuery,
) (*utils.ListResult[*models.Product], error) {
	ctx, span := p.tracer.Start(ctx, "mongoProductRepository.SearchProducts")
	span.SetAttributes(attribute2.String("SearchText", searchText))
	defer span.End()

	result, err := p.mongoGenericRepository.Search(ctx, searchText, listQuery)
	if err != nil {
		return nil, utils2.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"error in the paginate",
			),
		)
	}

	p.log.Infow(
		fmt.Sprintf(
			"products loaded for search term '%s'",
			searchText,
		),
		logger.Fields{"ProductsResult": result},
	)

	span.SetAttributes(attribute.Object("ProductsResult", result))

	return result, nil
}

// GetProductByID gets a product by id from the database.
func (p *mongoProductRepository) GetProductByID(
	ctx context.Context,
	uuid string,
) (*models.Product, error) {
	ctx, span := p.tracer.Start(ctx, "mongoProductRepository.GetProductByID")
	span.SetAttributes(attribute2.String("ID", uuid))
	defer span.End()

	id, err := uuid2.FromString(uuid)
	if err != nil {
		return nil, err
	}

	product, err := p.mongoGenericRepository.GetById(ctx, id)
	if err != nil {
		return nil, utils2.TraceStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"can't find the product with id %s into the database.",
					uuid,
				),
			),
		)
	}

	span.SetAttributes(attribute.Object("Product", product))

	p.log.Infow(
		fmt.Sprintf("product with id %s laoded", uuid),
		logger.Fields{"Product": product, "ID": uuid},
	)

	return product, nil
}

// GetProductByProductId gets a product by product id from the database.
func (p *mongoProductRepository) GetProductByProductId(
	ctx context.Context,
	uuid string,
) (*models.Product, error) {
	productId := uuid
	ctx, span := p.tracer.Start(
		ctx,
		"mongoProductRepository.GetProductByProductId",
	)
	span.SetAttributes(attribute2.String("ProductID", productId))
	defer span.End()

	product, err := p.mongoGenericRepository.FirstOrDefault(
		ctx,
		map[string]interface{}{"productId": uuid},
	)
	if err != nil {
		return nil, utils2.TraceStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"can't find the product with productId %s into the database.",
					uuid,
				),
			),
		)
	}

	span.SetAttributes(attribute.Object("Product", product))

	p.log.Infow(
		fmt.Sprintf(
			"product with productId %s laoded",
			productId,
		),
		logger.Fields{"Product": product, "ProductID": uuid},
	)

	return product, nil
}

// CreateProduct creates a product in the database.
func (p *mongoProductRepository) CreateProduct(
	ctx context.Context,
	product *models.Product,
) (*models.Product, error) {
	ctx, span := p.tracer.Start(ctx, "mongoProductRepository.CreateProduct")
	defer span.End()

	err := p.mongoGenericRepository.Add(ctx, product)
	if err != nil {
		return nil, utils2.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"error in the inserting product into the database.",
			),
		)
	}

	span.SetAttributes(attribute.Object("Product", product))

	p.log.Infow(
		fmt.Sprintf(
			"product with id '%s' created",
			product.ProductID,
		),
		logger.Fields{"Product": product, "ID": product.ProductID},
	)

	return product, nil
}

// UpdateProduct updates a product in the database.
func (p *mongoProductRepository) UpdateProduct(
	ctx context.Context,
	updateProduct *models.Product,
) (*models.Product, error) {
	ctx, span := p.tracer.Start(ctx, "mongoProductRepository.UpdateProduct")
	defer span.End()

	err := p.mongoGenericRepository.Update(ctx, updateProduct)
	// https://www.mongodb.com/docs/manual/reference/method/db.collection.findOneAndUpdate/
	if err != nil {
		return nil, utils2.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"error in updating product with id %s into the database.",
					updateProduct.ProductID,
				),
			),
		)
	}

	span.SetAttributes(attribute.Object("Product", updateProduct))
	p.log.Infow(
		fmt.Sprintf(
			"product with id '%s' updated",
			updateProduct.ProductID,
		),
		logger.Fields{"Product": updateProduct, "ID": updateProduct.ProductID},
	)

	return updateProduct, nil
}

// DeleteProductByID deletes a product by id from the database.
func (p *mongoProductRepository) DeleteProductByID(
	ctx context.Context,
	uuid string,
) error {
	ctx, span := p.tracer.Start(ctx, "mongoProductRepository.DeleteProductByID")
	span.SetAttributes(attribute2.String("ID", uuid))
	defer span.End()

	id, err := uuid2.FromString(uuid)
	if err != nil {
		return err
	}

	err = p.mongoGenericRepository.Delete(ctx, id)
	if err != nil {
		return utils2.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(err, fmt.Sprintf(
				"error in deleting product with id %s from the database.",
				uuid,
			)),
		)
	}

	p.log.Infow(
		fmt.Sprintf("product with id %s deleted", uuid),
		logger.Fields{"Product": uuid},
	)

	return nil
}
