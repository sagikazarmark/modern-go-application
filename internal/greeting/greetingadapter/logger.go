package greetingadapter

import (
	"github.com/InVisionApp/go-logger"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

// Logger wraps a go-kit logger and exposes it under a custom interface.
type Logger struct {
	logger log.Logger
}

// NewLogger returns a new Logger instance.
func NewLogger(logger log.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

// Debug logs a debug event.
func (l *Logger) Debug(msg ...interface{}) {
	l.logger.Debug(msg)
}

// Info logs an info event.
func (l *Logger) Info(msg ...interface{}) {
	l.logger.Info(msg)
}

// Warn logs a warning event.
func (l *Logger) Warn(msg ...interface{}) {
	l.logger.Warn(msg)
}

// Error logs an error event.
func (l *Logger) Error(msg ...interface{}) {
	l.logger.Error(msg)
}

// WithFields annotates a logger with some context.
func (l *Logger) WithFields(fields map[string]interface{}) greeting.Logger {
	return &Logger{logger: l.logger.WithFields(fields)}
}

// NewNopLogger returns a logger that doesn't do anything.
func NewNopLogger() *Logger {
	return NewLogger(log.NewNoop())
}
