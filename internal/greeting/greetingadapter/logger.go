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

// Debugf logs a debug event and optionally formats the message.
func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.logger.Debugf(msg, args...)
}

// Infof logs an info event and optionally formats the message.
func (l *Logger) Infof(msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

// Warnf logs a warning event and optionally formats the message.
func (l *Logger) Warnf(msg string, args ...interface{}) {
	l.logger.Warnf(msg, args...)
}

// Errorf logs an error event and optionally formats the message.
func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.logger.Errorf(msg, args...)
}

// WithFields annotates a logger with some context.
func (l *Logger) WithFields(fields map[string]interface{}) greeting.Logger {
	return &Logger{logger: l.logger.WithFields(fields)}
}

// NewNopLogger returns a logger that doesn't do anything.
func NewNopLogger() *Logger {
	return NewLogger(log.NewNoop())
}
