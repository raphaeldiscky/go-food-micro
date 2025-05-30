// Package problemdetails provides problem details.
package problemdetails

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"emperror.dev/errors"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

var Logger = defaultlogger.GetLogger()

// ContentTypeJSON is the content type for the JSON.
const (
	ContentTypeJSON = "application_exceptions/problem+json"
)

// ProblemDetailFunc is a function that returns a problem detail error.
type ProblemDetailFunc[E error] func(err E) ProblemDetailErr

// internalErrorMaps is a map of internal error maps.
var internalErrorMaps = map[reflect.Type]func(err error) ProblemDetailErr{}

// ProblemDetailErr ProblemDetail error interface.
type ProblemDetailErr interface {
	GetStatus() int
	SetStatus(status int) ProblemDetailErr
	GetTitle() string
	SetTitle(title string) ProblemDetailErr
	GetStackTrace() string
	SetStackTrace(stackTrace string) ProblemDetailErr
	GetDetail() string
	SetDetail(detail string) ProblemDetailErr
	GetType() string
	SetType(typ string) ProblemDetailErr
	Error() string
	ErrBody() error
}

// ProblemDetail error struct.
type problemDetail struct {
	Status     int       `json:"status,omitempty"`
	Title      string    `json:"title,omitempty"`
	Detail     string    `json:"detail,omitempty"`
	Type       string    `json:"type,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
	StackTrace string    `json:"stackTrace,omitempty"`
}

// ErrBody Error body.
func (p *problemDetail) ErrBody() error {
	return p
}

// Error  Error() interface method.
func (p *problemDetail) Error() string {
	return fmt.Sprintf(
		"Error Title: %s - Error Status: %d - Error Detail: %s",
		p.Title,
		p.Status,
		p.Detail,
	)
}

// GetStatus gets the status.
func (p *problemDetail) GetStatus() int {
	return p.Status
}

// SetStatus sets the status.
func (p *problemDetail) SetStatus(status int) ProblemDetailErr {
	p.Status = status

	return p
}

// GetTitle gets the title.
func (p *problemDetail) GetTitle() string {
	return p.Title
}

// SetTitle sets the title.
func (p *problemDetail) SetTitle(title string) ProblemDetailErr {
	p.Title = title

	return p
}

// GetType gets the type.
func (p *problemDetail) GetType() string {
	return p.Type
}

// SetType sets the type.
func (p *problemDetail) SetType(typ string) ProblemDetailErr {
	p.Type = typ

	return p
}

// GetDetail gets the detail.
func (p *problemDetail) GetDetail() string {
	return p.Detail
}

// SetDetail sets the detail.
func (p *problemDetail) SetDetail(detail string) ProblemDetailErr {
	p.Detail = detail

	return p
}

// GetStackTrace gets the stack trace.
func (p *problemDetail) GetStackTrace() string {
	return p.StackTrace
}

// SetStackTrace sets the stack trace.
func (p *problemDetail) SetStackTrace(stackTrace string) ProblemDetailErr {
	p.StackTrace = stackTrace

	return p
}

// NewProblemDetail New ProblemDetail Error.
func NewProblemDetail(
	status int,
	title string,
	detail string,
	stackTrace string,
) ProblemDetailErr {
	problemDetail := &problemDetail{
		Status:     status,
		Title:      title,
		Timestamp:  time.Now(),
		Detail:     detail,
		Type:       getDefaultType(status),
		StackTrace: stackTrace,
	}

	return problemDetail
}

// NewProblemDetailFromCode New ProblemDetail Error With Message.
func NewProblemDetailFromCode(status int, stackTrace string) ProblemDetailErr {
	return &problemDetail{
		Status:     status,
		Title:      http.StatusText(status),
		Timestamp:  time.Now(),
		Type:       getDefaultType(status),
		StackTrace: stackTrace,
	}
}

// NewProblemDetailFromCodeAndDetail New ProblemDetail Error With Message.
func NewProblemDetailFromCodeAndDetail(
	status int,
	detail string,
	stackTrace string,
) ProblemDetailErr {
	return &problemDetail{
		Status:     status,
		Title:      http.StatusText(status),
		Detail:     detail,
		Timestamp:  time.Now(),
		Type:       getDefaultType(status),
		StackTrace: stackTrace,
	}
}

// Map maps the problem detail.
func Map[E error](problem ProblemDetailFunc[E]) {
	errorType := typeMapper.GetGenericTypeByT[E]()
	if errorType.Kind() == reflect.Interface {
		types := typeMapper.TypesImplementedInterface[E]()
		for _, t := range types {
			internalErrorMaps[t] = func(err error) ProblemDetailErr {
				return problem(func() E {
					var target E
					_ = errors.As(err, &target)

					return target
				}())
			}
		}
	} else {
		internalErrorMaps[errorType] = func(err error) ProblemDetailErr {
			return problem(func() E {
				var target E
				_ = errors.As(err, &target)

				return target
			}())
		}
	}
}

// ResolveProblemDetail resolves the problem detail.
func ResolveProblemDetail(err error) ProblemDetailErr {
	resolvedErr := err
	for {
		_, ok := resolvedErr.(contracts.StackTracer)
		if ok {
			resolvedErr = errors.Unwrap(err)
		} else {
			break
		}
	}
	errorType := typeMapper.GetReflectType(resolvedErr)
	problem := internalErrorMaps[errorType]
	if problem != nil {
		return problem(resolvedErr)
	}

	return nil
}

// WriteTo writes the JSON Problem to an HTTP Response Writer.
func WriteTo(p ProblemDetailErr, w http.ResponseWriter) (int, error) {
	Logger.Error(p.Error())
	stackTrace := p.GetStackTrace()
	Logger.Infof("stackTrace: %s", stackTrace)

	writeHeaderTo(p, w)
	marshal, err := json.Marshal(&p)
	if err != nil {
		return 0, err
	}

	return w.Write(marshal)
}

// writeHeaderTo writes the header to the response writer.
func writeHeaderTo(p ProblemDetailErr, w http.ResponseWriter) {
	w.Header().Set("Content-Type", ContentTypeJSON)
	status := p.GetStatus()
	if status == 0 {
		status = http.StatusInternalServerError
	}

	w.WriteHeader(status)
}

// getDefaultType returns the default type.
func getDefaultType(statusCode int) string {
	return fmt.Sprintf("https://httpstatuses.io/%d", statusCode)
}
