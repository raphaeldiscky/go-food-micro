// Package customerrors provides custom errors.
package customerrors

import (
	"fmt"
	"io"
	"log"

	"emperror.dev/errors"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/contracts"
)

// customError is a struct that represents a custom error.
// https://klotzandrew.com/blog/error-handling-in-golang
// https://banzaicloud.com/blog/error-handling-go/
// https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
// https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
// https://github.com/go-stack/stack
// https://github.com/juju/errors
// https://github.com/emperror/errors
// https://github.com/pkg/errors/issues/75
type customError struct {
	statusCode int
	message    string
	error
}

// CustomError is a contract that represents a custom error.
type CustomError interface {
	error
	contracts.Wrapper
	contracts.Causer
	contracts.Formatter
	isCustomError()
	Status() int
	Message() string
}

// NewCustomError creates a new custom error.
func NewCustomError(err error, code int, message string) CustomError {
	m := &customError{
		statusCode: code,
		error:      err,
		message:    message,
	}

	return m
}

// isCustomError checks if the error is a custom error.
func (e *customError) isCustomError() {
}

// Error returns the error message.
func (e *customError) Error() string {
	if e.error != nil {
		return e.error.Error()
	}

	return e.message
}

// Message returns the error message.
func (e *customError) Message() string {
	return e.message
}

// Status returns the error status code.
func (e *customError) Status() int {
	return e.statusCode
}

// Cause returns the error cause.
func (e *customError) Cause() error {
	return e.error
}

// Unwrap returns the error unwrap.
func (e *customError) Unwrap() error {
	return e.error
}

// Format formats the error.
func (e *customError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// %s	error messages separated by a colon and a space (": ")
			// %q	double-quoted error messages separated by a colon and a space (": ")
			// %v	one error message per line
			// %+v	one error message per line and stack trace (if any)

			// if we have a call-stacked error, +v shows callstack for this error
			if _, err := fmt.Fprintf(s, "%+v", e.Cause()); err != nil {
				log.Printf("Error writing error string: %v", err)
			}
			// io.WriteString(s, e.message)
			return
		}

		fallthrough
	case 's', 'q':
		if _, err := io.WriteString(s, e.Error()); err != nil {
			log.Printf("Error writing error string: %v", err)
		}
	}
}

// GetCustomError gets the custom error.
func GetCustomError(err error) CustomError {
	if IsCustomError(err) {
		var internalErr CustomError
		errors.As(err, &internalErr)

		return internalErr
	}

	return nil
}

// IsCustomError checks if the error is a custom error.
func IsCustomError(err error) bool {
	var customErr CustomError

	var customError CustomError
	ok := errors.As(err, &customError)
	if ok {
		return true
	}

	// us, ok := errors.Cause(err).(ConflictError)
	if errors.As(err, &customErr) {
		return true
	}

	return false
}
