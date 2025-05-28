package v1

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
	"gorm.io/gorm"

	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	gormPostgres "github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/helpers/gormextensions"
	reflectionHelper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/reflectionhelper"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"

	datamodel "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/data/datamodels"
	dto "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/searchingproduct/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/models"
)

// searchProductsHandler is a struct that contains the search products handler.
type searchProductsHandler struct {
	fxparams.ProductHandlerParams
}

// NewSearchProductsHandler is a constructor for the searchProductsHandler.
func NewSearchProductsHandler(
	params fxparams.ProductHandlerParams,
) cqrs.RequestHandlerWithRegisterer[*SearchProducts, *dtos.SearchProductsResponseDto] {
	return &searchProductsHandler{
		ProductHandlerParams: params,
	}
}

// RegisterHandler is a method that registers the search products handler.
func (c *searchProductsHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*SearchProducts, *dtos.SearchProductsResponseDto](
		c,
	)
}

// Handle is a method that handles the search products query.
func (c *searchProductsHandler) Handle(
	ctx context.Context,
	query *SearchProducts,
) (*dtos.SearchProductsResponseDto, error) {
	dbQuery := c.prepareSearchDBQuery(query)

	products, err := gormPostgres.Paginate[*datamodel.ProductDataModel, *models.Product](
		ctx,
		query.ListQuery,
		dbQuery,
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in searching products in the repository",
		)
	}

	listResultDto, err := utils.ListResultToListResultDto[*dto.ProductDto](
		products,
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping ListResultToListResultDto",
		)
	}

	c.Log.Info("products fetched")

	return &dtos.SearchProductsResponseDto{Products: listResultDto}, nil
}

// prepareSearchDBQuery is a method that prepares the search db query.
func (c *searchProductsHandler) prepareSearchDBQuery(
	query *SearchProducts,
) *gorm.DB {
	fields := reflectionHelper.GetAllFields(
		typeMapper.GetGenericTypeByT[*datamodel.ProductDataModel](),
	)

	dbQuery := c.CatalogsDBContext.DB()

	for _, field := range fields {
		if field.Type.Kind() != reflect.String {
			continue
		}

		dbQuery = dbQuery.Or(
			fmt.Sprintf("%s LIKE ?", strcase.ToSnake(field.Name)),
			"%"+strings.ToLower(query.SearchText)+"%",
		)
	}

	return dbQuery
}
