package greetingworkeradapter

import (
	"github.com/goph/logur"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker"
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

// Debug logs an info event.
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	l.logger.Debug(msg, fields...)
}

// Info logs an info event.
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	l.logger.Info(msg, fields...)
}

// WithFields annotates a logger with some context.
func (l *Logger) WithFields(fields map[string]interface{}) greetingworker.Logger {
	return &Logger{logger: logur.WithFields(l.logger, fields)}
}

// NewNoopLogger returns a logger that discards all received log events.
func NewNoopLogger() *Logger {
	return NewLogger(logur.NewNoopLogger())
}
