// Package pipelines provides a set of functions for the pipelines.
// https://github.com/mehdihadeli/go-mediatr/blob/main/docs/pipelines.md
package pipelines

import (
	"context"

	"gorm.io/gorm"

	mediatr "github.com/mehdihadeli/go-mediatr"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/helpers/gormextensions"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// mediatorTransactionPipeline is a struct that contains the mediator transaction pipeline.
type mediatorTransactionPipeline struct {
	logger logger.Logger
	db     *gorm.DB
}

// NewMediatorTransactionPipeline creates a new mediator transaction pipeline.
func NewMediatorTransactionPipeline(
	l logger.Logger,
	db *gorm.DB,
) mediatr.PipelineBehavior {
	return &mediatorTransactionPipeline{
		logger: l,
		db:     db,
	}
}

// Handle handles the mediator transaction pipeline.
func (m *mediatorTransactionPipeline) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {
	requestName := typeMapper.GetSnakeTypeName(request)

	_, ok := request.(cqrs.TxRequest)
	if !ok {
		return next(ctx)
	}

	var result interface{}

	// https://gorm.io/docs/transactions.html#Transaction
	tx := m.db.WithContext(ctx).Begin()

	m.logger.Infof(
		"beginning database transaction for request `%s`",
		requestName,
	)

	gormContext := gormextensions.SetTxToContext(ctx, tx)
	ctx = gormContext

	defer func() {
		if r := recover(); r != nil {
			tx.WithContext(ctx).Rollback()

			if err, ok := r.(error); ok && err != nil {
				m.logger.Errorf(
					"panic tn the transaction, rolling back transaction with panic err: %+v",
					err,
				)
			} else {
				m.logger.Errorf("panic tn the transaction, rolling back transaction with panic message: %+v", r)
			}
		}
	}()

	middlewareResponse, err := next(ctx)
	result = middlewareResponse

	if err != nil {
		m.logger.Errorf(
			"rolling back transaction for request `%s`",
			requestName,
		)
		tx.WithContext(ctx).Rollback()

		return nil, err
	}

	m.logger.Infof("committing transaction for request `%s`", requestName)

	if err = tx.WithContext(ctx).Commit().Error; err != nil {
		m.logger.Errorf("transaction commit error: ", err)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}
