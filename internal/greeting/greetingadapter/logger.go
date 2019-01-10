package greetingadapter

import (
	"github.com/goph/logur"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

// Logger wraps a logur logger and exposes it under a custom interface.
type Logger struct {
	logger logur.Logger
}

// NewLogger returns a new Logger instance.
func NewLogger(logger logur.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

// Trace logs a trace event.
func (l *Logger) Trace(msg string, fields ...map[string]interface{}) {
	l.logger.Trace(msg, fields...)
}

// Debug logs a debug event.
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	l.logger.Debug(msg, fields...)
}

// Info logs an info event.
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	l.logger.Info(msg, fields...)
}

// Warn logs a warning event.
func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	l.logger.Warn(msg, fields...)
}

// Error logs an error event.
func (l *Logger) Error(msg string, fields ...map[string]interface{}) {
	l.logger.Error(msg, fields...)
}

// WithFields annotates a logger with some context.
func (l *Logger) WithFields(fields map[string]interface{}) greeting.Logger {
	return &Logger{logger: logur.WithFields(l.logger, fields)}
}

// NewNoopLogger returns a logger that discards all received log events.
func NewNoopLogger() *Logger {
	return NewLogger(logur.NewNoopLogger())
}
