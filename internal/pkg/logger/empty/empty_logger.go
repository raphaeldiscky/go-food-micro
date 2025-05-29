// Package empty provides an empty logger for the application.
package empty

import (
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/models"
)

var EmptyLogger logger.Logger = &emptyLogger{}

// emptyLogger is an empty logger.
type emptyLogger struct{}

// Configure configures the logger.
func (e emptyLogger) Configure(cfg func(internalLog interface{})) {
}

// Debug logs a debug message.
func (e emptyLogger) Debug(args ...interface{}) {
}

// Debugf logs a debug message with a format.
func (e emptyLogger) Debugf(template string, args ...interface{}) {
}

// Debugw logs a debug message with fields.
func (e emptyLogger) Debugw(msg string, fields logger.Fields) {
}

// LogType returns the log type.
func (e emptyLogger) LogType() models.LogType {
	return models.Zap
}

// Info logs an info message.
func (e emptyLogger) Info(args ...interface{}) {
}

// Infof logs an info message with a format.
func (e emptyLogger) Infof(template string, args ...interface{}) {
}

// Infow logs an info message with fields.
func (e emptyLogger) Infow(msg string, fields logger.Fields) {
}

// Warn logs a warning message.
func (e emptyLogger) Warn(args ...interface{}) {
}

// Warnf logs a warning message with a format.
func (e emptyLogger) Warnf(template string, args ...interface{}) {
}

// WarnMsg logs a warning message with an error.
func (e emptyLogger) WarnMsg(msg string, err error) {
}

// Error logs an error message.
func (e emptyLogger) Error(args ...interface{}) {
}

// Errorw logs an error message with fields.
func (e emptyLogger) Errorw(msg string, fields logger.Fields) {
}

// Errorf logs an error message with a format.
func (e emptyLogger) Errorf(template string, args ...interface{}) {
}

// Err logs an error message with an error.
func (e emptyLogger) Err(msg string, err error) {
}

// Fatal logs a fatal message.
func (e emptyLogger) Fatal(args ...interface{}) {
}

// Fatalf logs a fatal message with a format.
func (e emptyLogger) Fatalf(template string, args ...interface{}) {
}

// Printf logs a message with a format.
func (e emptyLogger) Printf(template string, args ...interface{}) {
}

// WithName logs a message with a name.
func (e emptyLogger) WithName(name string) {
}

// GrpcMiddlewareAccessLogger logs a grpc middleware access message.
func (e emptyLogger) GrpcMiddlewareAccessLogger(
	method string,
	time time.Duration,
	metaData map[string][]string,
	err error,
) {
}

// GrpcClientInterceptorLogger logs a grpc client interceptor message.
func (e emptyLogger) GrpcClientInterceptorLogger(
	method string,
	req interface{},
	reply interface{},
	time time.Duration,
	metaData map[string][]string,
	err error,
) {
}
