// Package logger provides a logger for the application.
package logger

import (
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/models"
)

// Fields is a map of fields.
type Fields map[string]interface{}

// Logger is a logger interface.
type Logger interface {
	Configure(cfg func(internalLog interface{}))
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, fields Fields)
	LogType() models.LogType
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, fields Fields)
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	WarnMsg(msg string, err error)
	Error(args ...interface{})
	Errorw(msg string, fields Fields)
	Errorf(template string, args ...interface{})
	Err(msg string, err error)
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Printf(template string, args ...interface{})
	WithName(name string)

	// GrpcMiddlewareAccessLogger logs a grpc middleware access message.
	GrpcMiddlewareAccessLogger(
		method string,
		time time.Duration,
		metaData map[string][]string,
		err error,
	)

	// GrpcClientInterceptorLogger logs a grpc client interceptor message.
	GrpcClientInterceptorLogger(
		method string,
		req interface{},
		reply interface{},
		time time.Duration,
		metaData map[string][]string,
		err error,
	)
}
