// Package scopes provides a set of functions for the scopes.
package scopes

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	uuid "github.com/satori/go.uuid"

	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
)

// AmountGreaterThan1000 returns the amount greater than 1000.
// https://gorm.io/docs/advanced_query.html#Scopes
// https://gorm.io/docs/scopes.html
// After scopes, we should have a runner function like Find, Update, Delete.
func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
	return db.Where("amount > ?", 1000)
}

// FilterAllItemsWithSoftDeleted returns soft-deleted and none soft-deleted items.
func FilterAllItemsWithSoftDeleted(db *gorm.DB) *gorm.DB {
	// https://gorm.io/docs/delete.html#Find-soft-deleted-records
	return db.Unscoped()
}

// SoftDeleted returns only soft-deleted items.
func SoftDeleted(db *gorm.DB) *gorm.DB {
	return db.Unscoped().Where("deleted_at IS NOT NULL")
}

// FilterByTitle filters the items by title.
func FilterByTitle(title string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("title = ?", title)
	}
}

// FilterByID filters the items by id.
func FilterByID(id uuid.UUID) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

// FilterPaginate filters the items by pagination.
func FilterPaginate[TDataModel any](
	ctx context.Context,
	listQuery *utils.ListQuery,
) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var totalRows int64

		dataModel := typeMapper.GenericInstanceByT[TDataModel]()
		// https://gorm.io/docs/advanced_query.html
		db.WithContext(ctx).Model(dataModel).Count(&totalRows)

		// generate where query
		query := db.WithContext(ctx).
			Model(dataModel).
			Offset(listQuery.GetOffset()).
			Limit(listQuery.GetLimit()).
			Order(listQuery.GetOrderBy())

		if listQuery.Filters != nil {
			for _, filter := range listQuery.Filters {
				column := filter.Field
				action := filter.Comparison
				value := filter.Value

				switch action {
				case "equals":
					whereQuery := fmt.Sprintf("%s = ?", column)
					query = query.Where(whereQuery, value)
				case "contains":
					whereQuery := fmt.Sprintf("%s LIKE ?", column)
					query = query.Where(whereQuery, "%"+value+"%")
				case "in":
					whereQuery := fmt.Sprintf("%s IN (?)", column)
					queryArray := strings.Split(value, ",")
					query = query.Where(whereQuery, queryArray)
				}
			}
		}

		return query
	}
}
