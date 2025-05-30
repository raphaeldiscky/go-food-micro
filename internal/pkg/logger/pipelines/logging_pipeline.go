// Package loggingpipelines provides a logging pipeline for the application.
package loggingpipelines

import (
	"context"
	"fmt"
	"time"

	mediatr "github.com/mehdihadeli/go-mediatr"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// requestLoggerPipeline is a request logger pipeline.
type requestLoggerPipeline struct {
	logger logger.Logger
}

// NewMediatorLoggingPipeline creates a new mediator logging pipeline.
func NewMediatorLoggingPipeline(l logger.Logger) mediatr.PipelineBehavior {
	return &requestLoggerPipeline{logger: l}
}

// Handle handles a request.
func (r *requestLoggerPipeline) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		r.logger.Infof("Request took %s", elapsed)
	}()

	requestName := typeMapper.GetNonePointerTypeName(request)

	r.logger.Infow(
		fmt.Sprintf("Handling request: '%s'", requestName),
		logger.Fields{"Request": request},
	)

	response, err := next(ctx)
	if err != nil {
		r.logger.Infof("Request failed with error: %v", err)

		return nil, err
	}

	responseName := typeMapper.GetNonePointerTypeName(response)

	r.logger.Infow(
		fmt.Sprintf(
			"Request handled successfully with response: '%s'",
			responseName,
		),
		logger.Fields{"Response": response},
	)

	return response, nil
}
