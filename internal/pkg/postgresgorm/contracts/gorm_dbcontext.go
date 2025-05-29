// Package contracts provides a set of contracts for the postgresgorm package.
package contracts

import (
	"context"

	"gorm.io/gorm"
)

// GormDBContext is a context that contains a gorm.DB.
type GormDBContext interface {
	WithTx(ctx context.Context) (GormDBContext, error)
	WithTxIfExists(ctx context.Context) GormDBContext
	RunInTx(ctx context.Context, action ActionFunc) error
	DB() *gorm.DB
}

// ActionFunc is a function that executes an action.
type ActionFunc func(ctx context.Context, gormContext GormDBContext) error
