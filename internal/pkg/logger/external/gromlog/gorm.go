// Package gromlog provides a logger for the application.
package gromlog

import (
	"context"
	"time"

	gormlogger "gorm.io/gorm/logger"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

// Ref: https://articles.wesionary.team/logging-interfaces-in-go-182c28be3d18
// implement gorm logger Interface.

// GormCustomLogger is a custom logger for the application.
type GormCustomLogger struct {
	logger.Logger
	gormlogger.Config
}

// NewGormCustomLogger creates a new custom logger for the application.
func NewGormCustomLogger(logger logger.Logger) *GormCustomLogger {
	// cfg, err := config.ProvideLogConfig()
	//
	// var logger logger.logger
	// if cfg.LogType == datamodels.Logrus && err != nil {
	//	logger = logrous.NewLogrusLogger(cfg, constants.Dev)
	// } else {
	//	if err != nil {
	//		cfg = &config.LogOptions{LogLevel: "info", LogType: datamodels.Zap}
	//	}
	//	logger = zap.NewZapLogger(cfg, constants.Dev)
	//}

	return &GormCustomLogger{
		Logger: logger,
		Config: gormlogger.Config{
			LogLevel: gormlogger.Info,
		},
	}
}

// LogMode sets log mode.
func (l *GormCustomLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level

	return &newlogger
}

// Info prints info messages.
func (l GormCustomLogger) Info(_ context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.Debugf(str, args...)
	}
}

// Warn prints warn messages.
func (l GormCustomLogger) Warn(_ context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.Warnf(str, args...)
	}
}

// Error prints error messages.
func (l GormCustomLogger) Error(_ context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.Errorf(str, args...)
	}
}

// Trace prints trace messages.
func (l GormCustomLogger) Trace(
	_ context.Context,
	begin time.Time,
	fc func() (string, int64),
	_ error,
) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	if l.LogLevel >= gormlogger.Info {
		sql, rows := fc()
		l.Debug("[", elapsed.Milliseconds(), " ms, ", rows, " rows] ", "sql -> ", sql)

		return
	}

	if l.LogLevel >= gormlogger.Warn {
		sql, rows := fc()
		l.Logger.Warn("[", elapsed.Milliseconds(), " ms, ", rows, " rows] ", "sql -> ", sql)

		return
	}

	if l.LogLevel >= gormlogger.Error {
		sql, rows := fc()
		l.Logger.Error("[", elapsed.Milliseconds(), " ms, ", rows, " rows] ", "sql -> ", sql)

		return
	}
}
