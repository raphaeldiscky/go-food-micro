// Package gormdbcontext provides a set of functions for the gormdbcontext package.
package gormdbcontext

import (
	"context"

	"gorm.io/gorm"

	defaultlogger "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/helpers/gormextensions"
)

// gormDBContext is a struct that contains a gorm.DB.
type gormDBContext struct {
	db *gorm.DB
}

// NewGormDBContext creates a new gormDBContext.
func NewGormDBContext(db *gorm.DB) contracts.GormDBContext {
	c := &gormDBContext{db: db}

	return c
}

// DB returns the gorm.DB.
func (c *gormDBContext) DB() *gorm.DB {
	return c.db
}

// WithTx creates a transactional DBContext with getting tx-gorm from the ctx. This will throw an error if the transaction does not exist.
func (c *gormDBContext) WithTx(
	ctx context.Context,
) (contracts.GormDBContext, error) {
	tx, err := gormextensions.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return NewGormDBContext(tx), nil
}

// WithTxIfExists creates a transactional DBContext with getting tx-gorm from the ctx. not throw an error if the transaction is not existing and returns an existing database.
func (c *gormDBContext) WithTxIfExists(
	ctx context.Context,
) contracts.GormDBContext {
	tx := gormextensions.GetTxFromContextIfExists(ctx)
	if tx == nil {
		return c
	}

	return NewGormDBContext(tx)
}

// RunInTx runs a transaction.
func (c *gormDBContext) RunInTx(
	ctx context.Context,
	action contracts.ActionFunc,
) error {
	// https://gorm.io/docs/transactions.html#Transaction
	tx := c.DB().WithContext(ctx).Begin()

	defaultlogger.GetLogger().Info("beginning database transaction")

	gormContext := gormextensions.SetTxToContext(ctx, tx)
	ctx = gormContext

	defer func() {
		if r := recover(); r != nil {
			tx.WithContext(ctx).Rollback()

			if err, ok := r.(error); ok && err != nil {
				defaultlogger.GetLogger().Errorf(
					"panic tn the transaction, rolling back transaction with panic err: %+v",
					err,
				)
			} else {
				defaultlogger.GetLogger().Errorf("panic tn the transaction, rolling back transaction with panic message: %+v", r)
			}
		}
	}()

	err := action(ctx, c)
	if err != nil {
		defaultlogger.GetLogger().Error("rolling back transaction")
		tx.WithContext(ctx).Rollback()

		return err
	}

	defaultlogger.GetLogger().Info("committing transaction")

	if err = tx.WithContext(ctx).Commit().Error; err != nil {
		defaultlogger.GetLogger().Errorf("transaction commit error: %+v", err)
	}

	return err
}
