// Package fxapp provides a module for the fxapp.
package fxapp

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

// FxErrorHandler is a fx error handler.
type FxErrorHandler struct {
	logger logger.Logger
}

// NewFxErrorHandler creates a new fx error handler.
func NewFxErrorHandler(logger logger.Logger) *FxErrorHandler {
	return &FxErrorHandler{logger: logger}
}

// HandleError handles an error.
func (h *FxErrorHandler) HandleError(e error) {
	h.logger.Error(e)
}
