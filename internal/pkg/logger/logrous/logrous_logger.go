// Package logrous provides a logger for the application.
package logrous

import (
	"os"
	"time"

	"github.com/nolleh/caption_json_formatter"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/constants"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	config2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/models"
)

// logrusLogger is a logrus logger.
type logrusLogger struct {
	level      string
	logger     *logrus.Logger
	logOptions *config2.LogOptions
}

// loggerLevelMap is a map of logger levels.
var loggerLevelMap = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
	"panic": logrus.PanicLevel,
	"fatal": logrus.FatalLevel,
}

// NewLogrusLogger creates a new logrus logger.
func NewLogrusLogger(
	cfg *config2.LogOptions,
	env environment.Environment,
) logger.Logger {
	logrusLogger := &logrusLogger{level: cfg.LogLevel, logOptions: cfg}
	logrusLogger.initLogger(env)

	return logrusLogger
}

// InitLogger Init logger.
func (l *logrusLogger) initLogger(env environment.Environment) {
	logLevel := l.GetLoggerLevel()

	// Create a new instance of the logger. You can have any number of instances.
	logrusLogger := logrus.New()

	logrusLogger.SetLevel(logLevel)

	// Output to stdout instead of the defaultLogger stderr
	// Can be any io.Writer, see below for File example
	logrusLogger.SetOutput(os.Stdout)

	if env.IsDevelopment() {
		logrusLogger.SetReportCaller(false)
		logrusLogger.SetFormatter(&logrus.TextFormatter{
			DisableColors: false,
			ForceColors:   true,
			FullTimestamp: true,
		})
	} else {
		logrusLogger.SetReportCaller(false)
		// https://github.com/nolleh/caption_json_formatter
		logrusLogger.SetFormatter(&caption_json_formatter.Formatter{PrettyPrint: true})
	}

	if l.logOptions.EnableTracing {
		// Instrument logrus.
		logrus.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		)))
	}

	l.logger = logrusLogger
}

// GetLoggerLevel gets the logger level.
func (l *logrusLogger) GetLoggerLevel() logrus.Level {
	level, exist := loggerLevelMap[l.level]
	if !exist {
		return logrus.DebugLevel
	}

	return level
}

// LogType returns the log type.
func (l *logrusLogger) LogType() models.LogType {
	return models.Logrus
}

// Configure configures the logger.
func (l *logrusLogger) Configure(cfg func(internalLog interface{})) {
	cfg(l.logger)
}

// Debug logs a debug message.
func (l *logrusLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Debugf logs a debug message with a format.
func (l *logrusLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

// Debugw logs a debug message with fields.
func (l *logrusLogger) Debugw(msg string, fields logger.Fields) {
	entry := l.mapToFields(fields)
	entry.Debug(msg)
}

// Info logs an info message.
func (l *logrusLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Infof logs an info message with a format.
func (l *logrusLogger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

// Infow logs an info message with fields.
func (l *logrusLogger) Infow(msg string, fields logger.Fields) {
	entry := l.mapToFields(fields)
	entry.Info(msg)
}

// Warn logs a warning message.
func (l *logrusLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Warnf logs a warning message with a format.
func (l *logrusLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

// WarnMsg logs a warning message with an error.
func (l *logrusLogger) WarnMsg(msg string, err error) {
	l.logger.Warn(msg, logrus.WithField("error", err.Error()))
}

// Error logs an error message.
func (l *logrusLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Errorw logs an error message with fields.
func (l *logrusLogger) Errorw(msg string, fields logger.Fields) {
	entry := l.mapToFields(fields)
	entry.Error(msg)
}

// Errorf logs an error message with a format.
func (l *logrusLogger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

// Err logs an error message with an error.
func (l *logrusLogger) Err(msg string, err error) {
	l.logger.Error(msg, logrus.WithField("error", err.Error()))
}

// Fatal logs a fatal message.
func (l *logrusLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Fatalf logs a fatal message with a format.
func (l *logrusLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

// Printf logs a message with a format.
func (l *logrusLogger) Printf(template string, args ...interface{}) {
	l.logger.Printf(template, args...)
}

// WithName logs a message with a name.
func (l *logrusLogger) WithName(name string) {
	l.logger.WithField(constants.NAME, name)
}

// GrpcMiddlewareAccessLogger logs a grpc middleware access message.
func (l *logrusLogger) GrpcMiddlewareAccessLogger(
	method string,
	time time.Duration,
	metaData map[string][]string,
	err error,
) {
	l.Info(
		constants.GRPC,
		logrus.WithField(constants.METHOD, method),
		logrus.WithField(constants.TIME, time),
		logrus.WithField(constants.METADATA, metaData),
		logrus.WithError(err),
	)
}

// GrpcClientInterceptorLogger logs a grpc client interceptor message.
func (l *logrusLogger) GrpcClientInterceptorLogger(
	method string,
	req interface{},
	reply interface{},
	time time.Duration,
	metaData map[string][]string,
	err error,
) {
	l.Info(
		constants.GRPC,
		logrus.WithField(constants.METHOD, method),
		logrus.WithField(constants.REQUEST, req),
		logrus.WithField(constants.REPLY, reply),
		logrus.WithField(constants.TIME, time),
		logrus.WithField(constants.METADATA, metaData),
		logrus.WithError(err),
	)
}

// mapToFields maps fields to logrus fields.
func (l *logrusLogger) mapToFields(
	fields map[string]interface{},
) *logrus.Entry {
	return l.logger.WithFields(logrus.Fields{"ss": 1})
}
