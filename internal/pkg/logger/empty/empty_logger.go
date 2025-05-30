// Package empty provides an empty logger for the application.
package empty

import (
	"log"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/models"
)

var EmptyLogger logger.Logger = &emptyLogger{}

// emptyLogger is an empty logger.
type emptyLogger struct{}

// Configure configures the logger.
func (e emptyLogger) Configure(_ func(internalLog interface{})) {
	log.Println("Configure")
}

// Debug logs a debug message.
func (e emptyLogger) Debug(args ...interface{}) {
	log.Println(args...)
}

// Debugf logs a debug message with a format.
func (e emptyLogger) Debugf(template string, args ...interface{}) {
	log.Printf(template, args...)
}

// Debugw logs a debug message with fields.
func (e emptyLogger) Debugw(msg string, fields logger.Fields) {
	log.Println(msg, fields)
}

// LogType returns the log type.
func (e emptyLogger) LogType() models.LogType {
	return models.Zap
}

// Info logs an info message.
func (e emptyLogger) Info(args ...interface{}) {
	log.Println(args...)
}

// Infof logs an info message with a format.
func (e emptyLogger) Infof(template string, args ...interface{}) {
	log.Printf(template, args...)
}

// Infow logs an info message with fields.
func (e emptyLogger) Infow(msg string, fields logger.Fields) {
	log.Println(msg, fields)
}

// Warn logs a warning message.
func (e emptyLogger) Warn(args ...interface{}) {
	log.Println(args...)
}

// Warnf logs a warning message with a format.
func (e emptyLogger) Warnf(template string, args ...interface{}) {
	log.Printf(template, args...)
}

// WarnMsg logs a warning message with an error.
func (e emptyLogger) WarnMsg(msg string, err error) {
	log.Println(msg, err)
}

// Error logs an error message.
func (e emptyLogger) Error(args ...interface{}) {
	log.Println(args...)
}

// Errorw logs an error message with fields.
func (e emptyLogger) Errorw(msg string, fields logger.Fields) {
	log.Println(msg, fields)
}

// Errorf logs an error message with a format.
func (e emptyLogger) Errorf(template string, args ...interface{}) {
	log.Printf(template, args...)
}

// Err logs an error message with an error.
func (e emptyLogger) Err(msg string, err error) {
	log.Println(msg, err)
}

// Fatal logs a fatal message.
func (e emptyLogger) Fatal(args ...interface{}) {
	log.Println(args...)
}

// Fatalf logs a fatal message with a format.
func (e emptyLogger) Fatalf(template string, args ...interface{}) {
	log.Printf(template, args...)
}

// Printf logs a message with a format.
func (e emptyLogger) Printf(template string, args ...interface{}) {
	log.Printf(template, args...)
}

// WithName logs a message with a name.
func (e emptyLogger) WithName(name string) {
	log.Println(name)
}

// GrpcMiddlewareAccessLogger logs a grpc middleware access message.
func (e emptyLogger) GrpcMiddlewareAccessLogger(
	method string,
	time time.Duration,
	metaData map[string][]string,
	err error,
) {
	log.Println(method, time, metaData, err)
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
	log.Println(method, req, reply, time, metaData, err)
}
