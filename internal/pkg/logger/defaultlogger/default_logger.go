// Package defaultlogger provides a default logger for the application.
package defaultLogger

import (
	"os"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/constants"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/logrous"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/models"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/zap"
)

var l logger.Logger

// initLogger initializes the logger.
func initLogger() {
	logType := os.Getenv("LogConfig_LogType")

	switch logType {
	case "Zap", "":
		l = zap.NewZapLogger(
			&config.LogOptions{LogType: models.Zap, CallerEnabled: false},
			constants.Dev,
		)
	case "Logrus":
		l = logrous.NewLogrusLogger(
			&config.LogOptions{LogType: models.Logrus, CallerEnabled: false},
			constants.Dev,
		)
	default:
	}
}

// GetLogger gets the logger.
func GetLogger() logger.Logger {
	if l == nil {
		initLogger()
	}

	return l
}
