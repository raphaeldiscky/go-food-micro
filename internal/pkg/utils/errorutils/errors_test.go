// Package errorutils provides a error utils.
package errorutils

import (
	"fmt"
	"testing"

	"emperror.dev/errors"
	"github.com/stretchr/testify/assert"
)

// TestStackTraceWithErrors tests the stack trace with errors.
func TestStackTraceWithErrors(_ *testing.T) {
	err := errors.WithMessage(errors.New("handling bad request"), "this is a bad-request")
	err = errors.WrapIf(err, "outer error message")

	res := ErrorsWithStack(err)
	fmt.Println(res)
}

// TestStackTrace tests the stack trace.
func TestStackTrace(_ *testing.T) {
	err := errors.WithMessage(errors.New("handling bad request"), "this is a bad-request")
	err = errors.WrapIf(err, "outer error message")

	res := StackTrace(err)
	fmt.Println(res)
}

// TestRootStackTrace tests the root stack trace.
func TestRootStackTrace(_ *testing.T) {
	err := errors.WithMessage(errors.New("handling bad request"), "this is a bad-request")
	err = errors.WrapIf(err, "outer error message")

	res := RootStackTrace(err)
	fmt.Println(res)
}

// TestAllLevelStackTrace tests the all level stack trace.
func TestAllLevelStackTrace(_ *testing.T) {
	err := errors.WrapIf(errors.New("handling bad request"), "this is a bad-request")
	err = errors.WrapIf(err, "outer error message")

	res := ErrorsWithStack(err)
	fmt.Println(res)
}

// TestErrorsWithoutStackTrace tests the errors without stack trace.
func TestErrorsWithoutStackTrace(t *testing.T) {
	err := errors.WrapIf(errors.New("handling bad request"), "this is a bad-request")
	err = errors.WrapIf(err, "outer error message")

	res := ErrorsWithoutStack(err, true)
	fmt.Println(res)
	assert.Contains(t, res, "outer error message\nthis is a bad-request\nhandling bad request")

	res = ErrorsWithoutStack(err, false)
	fmt.Println(res)
	assert.Contains(t, res, "outer error message: this is a bad-request: handling bad request")
}
