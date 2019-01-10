package greeting

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

	// WithFields annotates a logger with some context.
	WithFields(fields map[string]interface{}) Logger
}
