// Package data provides a module for the data.
package data

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/data/specification"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
)

// GenericRepositoryWithDataModel is a generic repository with data model.
type GenericRepositoryWithDataModel[TDataModel interface{}, TEntity interface{}] interface {
	Add(ctx context.Context, entity TEntity) error
	AddAll(ctx context.Context, entities []TEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (TEntity, error)
	GetByFilter(ctx context.Context, filters map[string]interface{}) ([]TEntity, error)
	GetByFuncFilter(ctx context.Context, filterFunc func(TEntity) bool) ([]TEntity, error)
	GetAll(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[TEntity], error)
	FirstOrDefault(ctx context.Context, filters map[string]interface{}) (TEntity, error)
	Search(
		ctx context.Context,
		searchTerm string,
		listQuery *utils.ListQuery,
	) (*utils.ListResult[TEntity], error)
	Update(ctx context.Context, entity TEntity) error
	UpdateAll(ctx context.Context, entities []TEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
	SkipTake(ctx context.Context, skip int, take int) ([]TEntity, error)
	Count(ctx context.Context) int64
	Find(ctx context.Context, specification specification.Specification) ([]TEntity, error)
}

// GenericRepository is a generic repository.
type GenericRepository[TEntity interface{}] interface {
	GenericRepositoryWithDataModel[TEntity, TEntity]
}
