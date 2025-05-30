// Package models provides models for the logger.
package models

// LogType is a log type.
type LogType int32

// LogType constants.
const (
	Zap    LogType = 0
	Logrus LogType = 1
)
