// Package errorutils provides a error utils.
package errorutils

import (
	"testing"

	"emperror.dev/errors"
	"github.com/stretchr/testify/assert"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
)

var Logger = defaultlogger.GetLogger()

// TestStackTraceWithErrors tests the stack trace with errors.
func TestStackTraceWithErrors(_ *testing.T) {
	err := errors.WithMessage(errors.New("handling bad request"), "this is a bad-request")
	err = errors.WrapIf(err, "outer error message")

	res := ErrorsWithStack(err)
	Logger.Info(res)
}

// TestStackTrace tests the stack trace.
func TestStackTrace(_ *testing.T) {
	err := errors.WithMessage(errors.New("handling bad request"), "this is a bad-request")
	err = errors.WrapIf(err, "outer error message")

	res := StackTrace(err)
	Logger.Info(res)
}

// TestRootStackTrace tests the root stack trace.
func TestRootStackTrace(_ *testing.T) {
	err := errors.WithMessage(errors.New("handling bad request"), "this is a bad-request")
	err = errors.WrapIf(err, "outer error message")

	res := RootStackTrace(err)
	Logger.Info(res)
}

// TestAllLevelStackTrace tests the all level stack trace.
func TestAllLevelStackTrace(_ *testing.T) {
	err := errors.WrapIf(errors.New("handling bad request"), "this is a bad-request")
	err = errors.WrapIf(err, "outer error message")

	res := ErrorsWithStack(err)
	Logger.Info(res)
}

// TestErrorsWithoutStackTrace tests the errors without stack trace.
func TestErrorsWithoutStackTrace(t *testing.T) {
	err := errors.WrapIf(errors.New("handling bad request"), "this is a bad-request")
	err = errors.WrapIf(err, "outer error message")

	res := ErrorsWithoutStack(err, true)
	Logger.Info(res)
	assert.Contains(t, res, "outer error message\nthis is a bad-request\nhandling bad request")

	res = ErrorsWithoutStack(err, false)
	Logger.Info(res)
	assert.Contains(t, res, "outer error message: this is a bad-request: handling bad request")
}
