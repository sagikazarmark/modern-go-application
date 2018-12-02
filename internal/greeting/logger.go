package greeting

// Logger is the fundamental interface for all log operations.
type Logger interface {
	// Debug logs a debug event.
	Debug(msg ...interface{})

	// Info logs an info event.
	Info(msg ...interface{})

	// Warn logs a warning event.
	Warn(msg ...interface{})

	// Error logs an error event.
	Error(msg ...interface{})

	// WithFields annotates a logger with some context.
	WithFields(fields map[string]interface{}) Logger
}
