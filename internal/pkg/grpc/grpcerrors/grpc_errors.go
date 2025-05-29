// Package grpcerrors provides custom grpc errors.
package grpcerrors

import (
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	defaultLogger "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
)

// grpcErr is a struct that represents a grpc error.
type grpcErr struct {
	Status     codes.Code `json:"status,omitempty"`
	Title      string     `json:"title,omitempty"`
	Detail     string     `json:"detail,omitempty"`
	Timestamp  time.Time  `json:"timestamp,omitempty"`
	StackTrace string     `json:"stackTrace,omitempty"`
}

// GrpcErr is an interface that represents a grpc error.
type GrpcErr interface {
	GetStatus() codes.Code
	SetStatus(status codes.Code) GrpcErr
	GetTitle() string
	SetTitle(title string) GrpcErr
	GetStackTrace() string
	SetStackTrace(stackTrace string) GrpcErr
	GetDetail() string
	SetDetail(detail string) GrpcErr
	Error() string
	ErrBody() error
	ToJSON() string
	ToGrpcResponseErr() error
}

// NewGrpcError is a function that creates a new grpc error.
func NewGrpcError(
	status codes.Code,
	title string,
	detail string,
	stackTrace string,
) GrpcErr {
	grpcErr := &grpcErr{
		Status:     status,
		Title:      title,
		Timestamp:  time.Now(),
		Detail:     detail,
		StackTrace: stackTrace,
	}

	return grpcErr
}

// ErrBody Error body.
func (p *grpcErr) ErrBody() error {
	return p
}

// Error  Error() interface method.
func (p *grpcErr) Error() string {
	return fmt.Sprintf(
		"Error Title: %s - Error Status: %d - Error Detail: %s",
		p.Title,
		p.Status,
		p.Detail,
	)
}

// GetStatus is a function that returns the status.
func (p *grpcErr) GetStatus() codes.Code {
	return p.Status
}

// SetStatus is a function that sets the status.
func (p *grpcErr) SetStatus(status codes.Code) GrpcErr {
	p.Status = status

	return p
}

// GetTitle is a function that returns the title.
func (p *grpcErr) GetTitle() string {
	return p.Title
}

// SetTitle is a function that sets the title.
func (p *grpcErr) SetTitle(title string) GrpcErr {
	p.Title = title

	return p
}

// GetDetail is a function that returns the detail.
func (p *grpcErr) GetDetail() string {
	return p.Detail
}

// SetDetail is a function that sets the detail.
func (p *grpcErr) SetDetail(detail string) GrpcErr {
	p.Detail = detail

	return p
}

// GetStackTrace is a function that returns the stack trace.
func (p *grpcErr) GetStackTrace() string {
	return p.StackTrace
}

// SetStackTrace is a function that sets the stack trace.
func (p *grpcErr) SetStackTrace(stackTrace string) GrpcErr {
	p.StackTrace = stackTrace

	return p
}

// ToGrpcResponseErr creates a gRPC error response to send grpc engine.
func (p *grpcErr) ToGrpcResponseErr() error {
	return status.Error(p.GetStatus(), p.ToJSON())
}

func (p *grpcErr) ToJSON() string {
	defaultLogger.GetLogger().Error(p.Error())
	stackTrace := p.GetStackTrace()
	fmt.Println(stackTrace)

	return string(p.json())
}

func (p *grpcErr) json() []byte {
	b, err := json.Marshal(&p)
	if err != nil {
		return []byte{}
	}

	return b
}
