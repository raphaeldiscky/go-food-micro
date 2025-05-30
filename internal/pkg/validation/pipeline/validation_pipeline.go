// Package pipeline provides a validation pipeline.
package pipeline

import (
	"context"

	mediatr "github.com/mehdihadeli/go-mediatr"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/validation"
)

// mediatorValidationPipeline is a struct that represents a mediator validation pipeline.
type mediatorValidationPipeline struct {
	logger logger.Logger
}

// NewMediatorValidationPipeline is a function that creates a new mediator validation pipeline.
func NewMediatorValidationPipeline(l logger.Logger) mediatr.PipelineBehavior {
	return &mediatorValidationPipeline{logger: l}
}

// Handle is a function that handles the request.
func (m mediatorValidationPipeline) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {
	v, ok := request.(validation.Validator)
	if ok {
		err := v.Validate()
		if err != nil {
			return nil, err
		}
	}

	return next(ctx)
}
