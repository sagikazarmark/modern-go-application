package greetingworker

// Logger is the fundamental interface for all log operations.
type Logger interface {
	// Info logs an info event.
	Info(msg ...interface{})

	// WithFields annotates a logger with some context.
	WithFields(fields map[string]interface{}) Logger
}

// LogFields is a shorthand for map[string]interface{} used in structured logging.
type LogFields map[string]interface{}
