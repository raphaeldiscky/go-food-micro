// Package repositories contains the postgres product repository.
package repositories

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/data"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/attribute"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/repository"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
	"gorm.io/gorm"

	utils2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/utils"
	goUuid "github.com/satori/go.uuid"
	attribute2 "go.opentelemetry.io/otel/attribute"

	data2 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/models"
)

// PostgresProductRepository is a struct that contains the postgres product repository.
type PostgresProductRepository struct {
	Log                   logger.Logger
	GormGenericRepository data.GenericRepository[*models.Product]
	Tracer                tracing.AppTracer
}

// NewPostgresProductRepository is a constructor for the PostgresProductRepository.
func NewPostgresProductRepository(
	log logger.Logger,
	db *gorm.DB,
	tracer tracing.AppTracer,
) data2.ProductRepository {
	gormRepository := repository.NewGenericGormRepository[*models.Product](db)

	return &PostgresProductRepository{
		Log:                   log,
		GormGenericRepository: gormRepository,
		Tracer:                tracer,
	}
}

// GetAllProducts is a method that gets all products.
func (p *PostgresProductRepository) GetAllProducts(
	ctx context.Context,
	listQuery *utils.ListQuery,
) (*utils.ListResult[*models.Product], error) {
	ctx, span := p.Tracer.Start(ctx, "postgresProductRepository.GetAllProducts")
	defer span.End()

	result, err := p.GormGenericRepository.GetAll(ctx, listQuery)
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

	p.Log.Infow(
		"products loaded",
		logger.Fields{"ProductsResult": result},
	)

	span.SetAttributes(attribute.Object("ProductsResult", result))

	return result, nil
}

// SearchProducts is a method that searches for products.
func (p *PostgresProductRepository) SearchProducts(
	ctx context.Context,
	searchText string,
	listQuery *utils.ListQuery,
) (*utils.ListResult[*models.Product], error) {
	ctx, span := p.Tracer.Start(ctx, "postgresProductRepository.SearchProducts")
	span.SetAttributes(attribute2.String("SearchText", searchText))
	defer span.End()

	result, err := p.GormGenericRepository.Search(ctx, searchText, listQuery)
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

	p.Log.Infow(
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
func (p *PostgresProductRepository) GetProductByID(
	ctx context.Context,
	uuid goUuid.UUID,
) (*models.Product, error) {
	ctx, span := p.Tracer.Start(ctx, "postgresProductRepository.GetProductByID")
	span.SetAttributes(attribute2.String("ID", uuid.String()))
	defer span.End()

	product, err := p.GormGenericRepository.GetByID(ctx, uuid)
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
	p.Log.Infow(
		fmt.Sprintf(
			"product with id %s laoded",
			uuid.String(),
		),
		logger.Fields{"Product": product, "ID": uuid},
	)

	return product, nil
}

// CreateProduct creates a new product in the database.
// It first checks if a product with the same ID already exists.
// If the product exists, it returns a Conflict error.
// If the product doesn't exist, it creates a new product and returns it.
func (p *PostgresProductRepository) CreateProduct(
	ctx context.Context,
	product *models.Product,
) (*models.Product, error) {
	ctx, span := p.Tracer.Start(ctx, "postgresProductRepository.CreateProduct")
	defer span.End()

	// Check if product already exists
	exists, err := p.GormGenericRepository.GetByID(ctx, product.ID)
	if err != nil {
		// If it's a not found error, we can proceed with creation
		if !customerrors.IsNotFoundError(err) {
			return nil, err
		}
	} else if exists != nil {
		return nil, customerrors.NewConflictError(
			fmt.Sprintf("product with id '%s' already exists", product.ID),
		)
	}

	err = p.GormGenericRepository.Add(ctx, product)
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
	p.Log.Infow(
		fmt.Sprintf(
			"product with id '%s' created",
			product.ID,
		),
		logger.Fields{"Product": product, "ID": product.ID},
	)

	return product, nil
}

// UpdateProduct updates an existing product in the database.
// It updates all fields of the product and returns the updated product.
// If the product doesn't exist, it returns an error.
func (p *PostgresProductRepository) UpdateProduct(
	ctx context.Context,
	updateProduct *models.Product,
) (*models.Product, error) {
	ctx, span := p.Tracer.Start(ctx, "postgresProductRepository.UpdateProduct")
	defer span.End()

	err := p.GormGenericRepository.Update(ctx, updateProduct)
	err = utils2.TraceStatusFromSpan(
		span,
		errors.WrapIf(
			err,
			fmt.Sprintf(
				"error in updating product with id %s into the database.",
				updateProduct.ID,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Object("Product", updateProduct))
	p.Log.Infow(
		fmt.Sprintf(
			"product with id '%s' updated",
			updateProduct.ID,
		),
		logger.Fields{
			"Product": updateProduct,
			"ID":      updateProduct.ID,
		},
	)

	return updateProduct, nil
}

// DeleteProductByID deletes a product from the database by its ID.
// It first checks if the product exists.
// If the product doesn't exist, it returns a NotFound error.
// If the product exists, it deletes it and returns nil.
func (p *PostgresProductRepository) DeleteProductByID(
	ctx context.Context,
	uuid goUuid.UUID,
) error {
	ctx, span := p.Tracer.Start(ctx, "postgresProductRepository.DeleteProductByID")
	span.SetAttributes(attribute2.String("ID", uuid.String()))
	defer span.End()

	// Check if product exists before deleting
	exists, err := p.GormGenericRepository.GetByID(ctx, uuid)
	if err != nil || exists == nil {
		return customerrors.NewNotFoundError(
			fmt.Sprintf("product with id `%s` not found in the database", uuid),
		)
	}

	err = p.GormGenericRepository.Delete(ctx, uuid)
	err = utils2.TraceStatusFromSpan(span, errors.WrapIf(err, fmt.Sprintf(
		"error in deleting product with id `%s` from the database",
		uuid,
	)))
	if err != nil {
		return err
	}

	p.Log.Infow(
		fmt.Sprintf(
			"product with id `%s` deleted",
			uuid,
		),
		logger.Fields{"Product": uuid},
	)

	return nil
}
