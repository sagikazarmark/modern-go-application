package common

import (
	"context"
)

// Logger is the fundamental interface for all log operations.
type Logger interface {
	// Trace logs a debug event.
	Trace(msg string, fields ...map[string]interface{})

	// Debug logs a debug event.
	Debug(msg string, fields ...map[string]interface{})

	// Info logs an info event.
	Info(msg string, fields ...map[string]interface{})

	// Warn logs a warning event.
	Warn(msg string, fields ...map[string]interface{})

	// Error logs an error event.
	Error(msg string, fields ...map[string]interface{})

	// WithFields annotates a logger with key-value pairs.
	WithFields(fields map[string]interface{}) Logger

	// WithContext annotates a logger with a context.
	WithContext(ctx context.Context) Logger
}

type noopLogger struct{}

// NewNoopLogger returns a logger that discards every log event.
func NewNoopLogger() Logger { return noopLogger{} }

func (noopLogger) Trace(msg string, fields ...map[string]interface{}) {}
func (noopLogger) Debug(msg string, fields ...map[string]interface{}) {}
func (noopLogger) Info(msg string, fields ...map[string]interface{})  {}
func (noopLogger) Warn(msg string, fields ...map[string]interface{})  {}
func (noopLogger) Error(msg string, fields ...map[string]interface{}) {}
func (n noopLogger) WithFields(fields map[string]interface{}) Logger  { return n }
func (n noopLogger) WithContext(ctx context.Context) Logger           { return n }
