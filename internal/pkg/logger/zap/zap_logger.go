// Package zap provides a logger for the application.
package zap

import (
	"os"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/constants"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	config2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/models"
)

const (
	serviceKey = "[SERVICE]"
	timeKey    = "[TIME]"
	levelKey   = "[LEVEL]"
	callerKey  = "[CALLER]"
	lineKey    = "[LINE]"
	messageKey = "[MESSAGE]"
)

// zapLogger is a zap logger.
type zapLogger struct {
	level       string
	sugarLogger *zap.SugaredLogger
	logger      *zap.Logger
	logOptions  *config2.LogOptions
}

// ZapLogger is a zap logger.
type ZapLogger interface {
	logger.Logger
	InternalLogger() *zap.Logger
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Sync() error
}

// For mapping config logger.
var loggerLevelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"panic": zapcore.PanicLevel,
	"fatal": zapcore.FatalLevel,
}

// NewZapLogger create new zap logger.
func NewZapLogger(
	cfg *config2.LogOptions,
	env environment.Environment,
) ZapLogger {
	zapLogger := &zapLogger{level: cfg.LogLevel, logOptions: cfg}
	zapLogger.initLogger(env)

	return zapLogger
}

// InternalLogger returns the internal logger.
func (l *zapLogger) InternalLogger() *zap.Logger {
	return l.logger
}

// getLoggerLevel gets the logger level.
func (l *zapLogger) getLoggerLevel() zapcore.Level {
	level, exist := loggerLevelMap[l.level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

// InitLogger Init logger.
func (l *zapLogger) initLogger(env environment.Environment) {
	logLevel := l.getLoggerLevel()

	logWriter := zapcore.AddSync(os.Stdout)

	var encoderCfg zapcore.EncoderConfig
	var encoder zapcore.Encoder

	if env.IsProduction() {
		encoderCfg = zap.NewProductionEncoderConfig()
		encoderCfg.NameKey = serviceKey
		encoderCfg.TimeKey = timeKey
		encoderCfg.LevelKey = levelKey
		encoderCfg.FunctionKey = callerKey
		encoderCfg.CallerKey = lineKey
		encoderCfg.MessageKey = messageKey
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
		encoderCfg.EncodeName = zapcore.FullNameEncoder
		encoderCfg.EncodeDuration = zapcore.StringDurationEncoder
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
		encoderCfg.NameKey = serviceKey
		encoderCfg.TimeKey = timeKey
		encoderCfg.LevelKey = levelKey
		encoderCfg.FunctionKey = callerKey
		encoderCfg.CallerKey = lineKey
		encoderCfg.MessageKey = messageKey
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeName = zapcore.FullNameEncoder
		encoderCfg.EncodeDuration = zapcore.StringDurationEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderCfg.EncodeCaller = zapcore.FullCallerEncoder
		encoderCfg.ConsoleSeparator = " | "
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))

	var options []zap.Option

	if l.logOptions.CallerEnabled {
		options = append(options, zap.AddCaller())
		options = append(options, zap.AddCallerSkip(1))
	}

	logger := zap.New(core, options...)

	if l.logOptions.EnableTracing {
		// add logs as events to tracing
		logger = otelzap.New(logger).Logger
	}

	l.logger = logger
	l.sugarLogger = logger.Sugar()
}

// Configure configures the logger.
func (l *zapLogger) Configure(cfg func(internalLog interface{})) {
	cfg(l.logger)
}

// LogType returns the log type.
func (l *zapLogger) LogType() models.LogType {
	return models.Zap
}

// WithName add logger microservice name.
func (l *zapLogger) WithName(name string) {
	l.logger = l.logger.Named(name)
	l.sugarLogger = l.sugarLogger.Named(name)
}

// Debug uses fmt.Sprint to construct and log a message.
func (l *zapLogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (l *zapLogger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

// Debugw logs a message with fields.
func (l *zapLogger) Debugw(msg string, fields logger.Fields) {
	zapFields := mapToZapFields(fields)
	l.logger.Debug(msg, zapFields...)
}

// Info uses fmt.Sprint to construct and log a message.
func (l *zapLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *zapLogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

// Infow logs a message with some additional context.
func (l *zapLogger) Infow(msg string, fields logger.Fields) {
	zapFields := mapToZapFields(fields)
	l.logger.Info(msg, zapFields...)
}

// Printf uses fmt.Sprintf to log a templated message.
func (l *zapLogger) Printf(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *zapLogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

// WarnMsg log error message with warn level.
func (l *zapLogger) WarnMsg(msg string, err error) {
	l.logger.Warn(msg, zap.String("error", err.Error()))
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *zapLogger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

// Error uses fmt.Sprint to construct and log a message.
func (l *zapLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

// Errorw logs a message with some additional context.
func (l *zapLogger) Errorw(msg string, fields logger.Fields) {
	zapFields := mapToZapFields(fields)
	l.logger.Error(msg, zapFields...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *zapLogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

// Err uses error to log a message.
func (l *zapLogger) Err(msg string, err error) {
	l.logger.Error(msg, zap.Error(err))
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the logger then panics. (See DPanicLevel for details.)
func (l *zapLogger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the logger then panics. (See DPanicLevel for details.)
func (l *zapLogger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *zapLogger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (l *zapLogger) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *zapLogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *zapLogger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}

// Sync flushes any buffered log entries.
func (l *zapLogger) Sync() error {
	go func() {
		err := l.logger.Sync()
		if err != nil {
			l.logger.Error("error while syncing", zap.Error(err))
		}
	}()

	return l.sugarLogger.Sync()
}

// GrpcMiddlewareAccessLogger logs a grpc middleware access message.
func (l *zapLogger) GrpcMiddlewareAccessLogger(
	method string,
	t time.Duration,
	metaData map[string][]string,
	err error,
) {
	l.Info(
		constants.GRPC,
		zap.String(constants.METHOD, method),
		zap.Duration(constants.TIME, t),
		zap.Any(constants.METADATA, metaData),
		zap.Error(err),
	)
}

// GrpcClientInterceptorLogger logs a grpc client interceptor message.
func (l *zapLogger) GrpcClientInterceptorLogger(
	method string,
	req, reply interface{},
	time time.Duration,
	metaData map[string][]string,
	err error,
) {
	l.Info(
		constants.GRPC,
		zap.String(constants.METHOD, method),
		zap.Any(constants.REQUEST, req),
		zap.Any(constants.REPLY, reply),
		zap.Duration(constants.TIME, time),
		zap.Any(constants.METADATA, metaData),
		zap.Error(err),
	)
}

// mapToZapFields maps data to zap fields.
func mapToZapFields(data map[string]interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(data))

	for key, value := range data {
		field := zap.Field{
			Key:       key,
			Type:      getFieldType(value),
			Interface: value,
		}
		fields = append(fields, field)
	}

	return fields
}

func getFieldType(value interface{}) zapcore.FieldType {
	switch value.(type) {
	case string:
		return zapcore.StringType
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return zapcore.Int64Type
	case bool:
		return zapcore.BoolType
	case float32, float64:
		return zapcore.Float64Type
	case error:
		return zapcore.ErrorType
	default:
		// uses reflection to serialize arbitrary objects, so it can be slow and allocation-heavy.
		return zapcore.ReflectType
	}
}
