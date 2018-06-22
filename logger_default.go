package amqpx

import (
	"log"
	"os"
)

const (
	defaultLoggerPrefix = "[amqpx] "
)

// defaultLogger is the default Logger implementation using log.Logger.
type defaultLogger struct {
	level LoggerLevel
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
}

// newDefaultLogger returns a new defaultLogger instance.
func newDefaultLogger(level LoggerLevel) Logger {
	if level == LoggerLevelDisabled {
		return &noopLogger{}
	}

	flags := log.Ldate | log.Ltime

	return &defaultLogger{
		level: level,
		debug: log.New(os.Stdout, defaultLoggerPrefix+"DEBUG: ", flags),
		info:  log.New(os.Stdout, defaultLoggerPrefix+"INFO: ", flags),
		warn:  log.New(os.Stdout, defaultLoggerPrefix+"WARN: ", flags),
		err:   log.New(os.Stderr, defaultLoggerPrefix+"ERROR: ", flags),
	}
}

// Debug implements Logger interface.
func (l defaultLogger) Debug(args ...interface{}) {
	if l.level < LoggerLevelDebug {
		return
	}

	l.debug.Println(args...)
}

// Info implements Logger interface.
func (l defaultLogger) Info(args ...interface{}) {
	if l.level < LoggerLevelInfo {
		return
	}

	l.info.Println(args...)
}

// Warn implements Logger interface.
func (l defaultLogger) Warn(args ...interface{}) {
	if l.level < LoggerLevelWarn {
		return
	}

	l.warn.Println(args...)
}

// Error implements Logger interface.
func (l defaultLogger) Error(args ...interface{}) {
	if l.level < LoggerLevelError {
		return
	}

	l.err.Println(args...)
}

var _ Logger = (*defaultLogger)(nil)
