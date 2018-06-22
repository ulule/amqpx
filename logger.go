package amqpx

import (
	"strings"
)

// LoggerLevel represents a logging level.
type LoggerLevel uint8

// Logger levels
const (
	LoggerLevelDisabled LoggerLevel = iota
	LoggerLevelError
	LoggerLevelWarn
	LoggerLevelInfo
	LoggerLevelDebug
)

// Logger level strings
const (
	LoggerLevelDisabledStr = ""
	LoggerLevelErrorStr    = "error"
	LoggerLevelWarnStr     = "warn"
	LoggerLevelInfoStr     = "info"
	LoggerLevelDebugStr    = "debug"
)

// LoggerLevelFromString is a convenient method to returns a LoggerLevel from a string.
func LoggerLevelFromString(level string) LoggerLevel {
	level = strings.TrimSpace(strings.ToLower(level))
	switch level {
	case LoggerLevelDebugStr:
		return LoggerLevelDebug
	case LoggerLevelInfoStr:
		return LoggerLevelInfo
	case LoggerLevelWarnStr:
		return LoggerLevelWarn
	case LoggerLevelErrorStr:
		return LoggerLevelError
	default:
		return LoggerLevelDisabled
	}
}

// Logger describes a simple logger.
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}
