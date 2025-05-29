// Package problemdetails provides problem details.
package problemdetails

import (
	"net/http"
	"testing"

	"emperror.dev/errors"
	"github.com/stretchr/testify/assert"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
)

// TestDomainErr tests the domain error.
func TestDomainErr(t *testing.T) {
	domainErr := NewDomainProblemDetail(
		http.StatusBadRequest,
		"Order with id '1' already completed",
		"stack",
	)

	assert.Equal(t, "Order with id '1' already completed", domainErr.GetDetail())
	assert.Equal(t, "Domain Model Error", domainErr.GetTitle())
	assert.Equal(t, "stack", domainErr.GetStackTrace())
	assert.Equal(t, "https://httpstatuses.io/400", domainErr.GetType())
	assert.Equal(t, 400, domainErr.GetStatus())
}

// TestApplicationErr tests the application error.
func TestApplicationErr(t *testing.T) {
	applicationErr := NewApplicationProblemDetail(
		http.StatusBadRequest,
		"application_exceptions error",
		"stack",
	)

	assert.Equal(t, "application_exceptions error", applicationErr.GetDetail())
	assert.Equal(t, "Application Service Error", applicationErr.GetTitle())
	assert.Equal(t, "stack", applicationErr.GetStackTrace())
	assert.Equal(t, "https://httpstatuses.io/400", applicationErr.GetType())
	assert.Equal(t, 400, applicationErr.GetStatus())
}

// TestBadRequestErr tests the bad request error.
func TestBadRequestErr(t *testing.T) {
	badRequestError := NewBadRequestProblemDetail("bad-request error", "stack")

	assert.Equal(t, "bad-request error", badRequestError.GetDetail())
	assert.Equal(t, "Bad Request", badRequestError.GetTitle())
	assert.Equal(t, "stack", badRequestError.GetStackTrace())
	assert.Equal(t, "https://httpstatuses.io/400", badRequestError.GetType())
	assert.Equal(t, 400, badRequestError.GetStatus())
}

// TestParseError tests the parse error.
func TestParseError(t *testing.T) {
	// Bad-Request ProblemDetail
	badRequestError := errors.WrapIf(
		customErrors.NewBadRequestError("bad-request error"),
		"bad request error",
	)
	badRequestPrb := ParseError(badRequestError)
	assert.NotNil(t, badRequestPrb)
	assert.Equal(t, badRequestPrb.GetStatus(), 400)

	// NotFound ProblemDetail
	notFoundError := customErrors.NewNotFoundError("notfound error")
	notfoundPrb := ParseError(notFoundError)
	assert.NotNil(t, notFoundError)
	assert.Equal(t, notfoundPrb.GetStatus(), 404)
}

// TestMap tests the map.
func TestMap(t *testing.T) {
	Map[customErrors.BadRequestError](func(err customErrors.BadRequestError) ProblemDetailErr {
		return NewBadRequestProblemDetail(err.Message(), err.Error())
	})
	s := ResolveProblemDetail(customErrors.NewBadRequestError(""))
	assert.NotNil(t, s)
	assert.IsType(t, (*ProblemDetailErr)(nil), s)
}
