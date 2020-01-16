package common

import (
	"context"
)

// Logger is the fundamental interface for all log operations.
type Logger interface {
	// Trace logs a trace event.
	Trace(msg string, fields ...map[string]interface{})

	// Debug logs a debug event.
	Debug(msg string, fields ...map[string]interface{})

	// Info logs an info event.
	Info(msg string, fields ...map[string]interface{})

	// Warn logs a warning event.
	Warn(msg string, fields ...map[string]interface{})

	// Error logs an error event.
	Error(msg string, fields ...map[string]interface{})

	// TraceContext logs a trace event with a context.
	TraceContext(ctx context.Context, msg string, fields ...map[string]interface{})

	// DebugContext logs a debug event with a context.
	DebugContext(ctx context.Context, msg string, fields ...map[string]interface{})

	// InfoContext logs an info event with a context.
	InfoContext(ctx context.Context, msg string, fields ...map[string]interface{})

	// WarnContext logs a warning event with a context.
	WarnContext(ctx context.Context, msg string, fields ...map[string]interface{})

	// ErrorContext logs an error event with a context.
	ErrorContext(ctx context.Context, msg string, fields ...map[string]interface{})

	// WithFields annotates a logger with key-value pairs.
	WithFields(fields map[string]interface{}) Logger

	// WithContext annotates a logger with a context.
	WithContext(ctx context.Context) Logger
}

// NoopLogger is a logger that discards every log event.
type NoopLogger struct{}

func (NoopLogger) Trace(_ string, _ ...map[string]interface{}) {}
func (NoopLogger) Debug(_ string, _ ...map[string]interface{}) {}
func (NoopLogger) Info(_ string, _ ...map[string]interface{})  {}
func (NoopLogger) Warn(_ string, _ ...map[string]interface{})  {}
func (NoopLogger) Error(_ string, _ ...map[string]interface{}) {}

func (NoopLogger) TraceContext(_ context.Context, _ string, _ ...map[string]interface{}) {}
func (NoopLogger) DebugContext(_ context.Context, _ string, _ ...map[string]interface{}) {}
func (NoopLogger) InfoContext(_ context.Context, _ string, _ ...map[string]interface{})  {}
func (NoopLogger) WarnContext(_ context.Context, _ string, _ ...map[string]interface{})  {}
func (NoopLogger) ErrorContext(_ context.Context, _ string, _ ...map[string]interface{}) {}

func (n NoopLogger) WithFields(_ map[string]interface{}) Logger { return n }
func (n NoopLogger) WithContext(_ context.Context) Logger       { return n }
