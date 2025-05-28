package repositories

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/data"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/attribute"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/repository"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
	"gorm.io/gorm"

	utils2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/utils"
	uuid "github.com/satori/go.uuid"
	attribute2 "go.opentelemetry.io/otel/attribute"

	data2 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/models"
)

// postgresProductRepository is a struct that contains the postgres product repository.
type postgresProductRepository struct {
	log                   logger.Logger
	gormGenericRepository data.GenericRepository[*models.Product]
	tracer                tracing.AppTracer
}

// NewPostgresProductRepository is a constructor for the postgresProductRepository.
func NewPostgresProductRepository(
	log logger.Logger,
	db *gorm.DB,
	tracer tracing.AppTracer,
) data2.ProductRepository {
	gormRepository := repository.NewGenericGormRepository[*models.Product](db)

	return &postgresProductRepository{
		log:                   log,
		gormGenericRepository: gormRepository,
		tracer:                tracer,
	}
}

// GetAllProducts is a method that gets all products.
func (p *postgresProductRepository) GetAllProducts(
	ctx context.Context,
	listQuery *utils.ListQuery,
) (*utils.ListResult[*models.Product], error) {
	ctx, span := p.tracer.Start(ctx, "postgresProductRepository.GetAllProducts")
	defer span.End()

	result, err := p.gormGenericRepository.GetAll(ctx, listQuery)
	err = utils2.TraceStatusFromContext(
		ctx,
		errors.WrapIf(
			err,
			"error in the paginate",
		),
	)
	if err != nil {
		return nil, err
	}

	p.log.Infow(
		"products loaded",
		logger.Fields{"ProductsResult": result},
	)

	span.SetAttributes(attribute.Object("ProductsResult", result))

	return result, nil
}

// SearchProducts is a method that searches for products.
func (p *postgresProductRepository) SearchProducts(
	ctx context.Context,
	searchText string,
	listQuery *utils.ListQuery,
) (*utils.ListResult[*models.Product], error) {
	ctx, span := p.tracer.Start(ctx, "postgresProductRepository.SearchProducts")
	span.SetAttributes(attribute2.String("SearchText", searchText))
	defer span.End()

	result, err := p.gormGenericRepository.Search(ctx, searchText, listQuery)
	err = utils2.TraceStatusFromContext(
		ctx,
		errors.WrapIf(
			err,
			"error in the paginate",
		),
	)
	if err != nil {
		return nil, err
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

// GetProductByID is a method that gets a product by id.
func (p *postgresProductRepository) GetProductByID(
	ctx context.Context,
	uuid uuid.UUID,
) (*models.Product, error) {
	ctx, span := p.tracer.Start(ctx, "postgresProductRepository.GetProductByID")
	span.SetAttributes(attribute2.String("Id", uuid.String()))
	defer span.End()

	product, err := p.gormGenericRepository.GetById(ctx, uuid)
	err = utils2.TraceStatusFromSpan(
		span,
		errors.WrapIf(
			err,
			fmt.Sprintf(
				"can't find the product with id %s into the database.",
				uuid,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Object("Product", product))
	p.log.Infow(
		fmt.Sprintf(
			"product with id %s laoded",
			uuid.String(),
		),
		logger.Fields{"Product": product, "Id": uuid},
	)

	return product, nil
}

func (p *postgresProductRepository) CreateProduct(
	ctx context.Context,
	product *models.Product,
) (*models.Product, error) {
	ctx, span := p.tracer.Start(ctx, "postgresProductRepository.CreateProduct")
	defer span.End()

	err := p.gormGenericRepository.Add(ctx, product)
	err = utils2.TraceStatusFromSpan(
		span,
		errors.WrapIf(
			err,
			"error in the inserting product into the database.",
		),
	)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Object("Product", product))
	p.log.Infow(
		fmt.Sprintf(
			"product with id '%s' created",
			product.Id,
		),
		logger.Fields{"Product": product, "Id": product.Id},
	)

	return product, nil
}

func (p *postgresProductRepository) UpdateProduct(
	ctx context.Context,
	updateProduct *models.Product,
) (*models.Product, error) {
	ctx, span := p.tracer.Start(ctx, "postgresProductRepository.UpdateProduct")
	defer span.End()

	err := p.gormGenericRepository.Update(ctx, updateProduct)
	err = utils2.TraceStatusFromSpan(
		span,
		errors.WrapIf(
			err,
			fmt.Sprintf(
				"error in updating product with id %s into the database.",
				updateProduct.Id,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Object("Product", updateProduct))
	p.log.Infow(
		fmt.Sprintf(
			"product with id '%s' updated",
			updateProduct.Id,
		),
		logger.Fields{
			"Product": updateProduct,
			"Id":      updateProduct.Id,
		},
	)

	return updateProduct, nil
}

func (p *postgresProductRepository) DeleteProductByID(
	ctx context.Context,
	uuid uuid.UUID,
) error {
	ctx, span := p.tracer.Start(ctx, "postgresProductRepository.UpdateProduct")
	span.SetAttributes(attribute2.String("Id", uuid.String()))
	defer span.End()

	err := p.gormGenericRepository.Delete(ctx, uuid)
	err = utils2.TraceStatusFromSpan(span, errors.WrapIf(err, fmt.Sprintf(
		"error in the deleting product with id %s into the database.",
		uuid,
	)))
	if err != nil {
		return err
	}

	p.log.Infow(
		fmt.Sprintf(
			"product with id %s deleted",
			uuid,
		),
		logger.Fields{"Product": uuid},
	)

	return nil
}
