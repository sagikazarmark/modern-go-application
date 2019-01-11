package runner

// Logger is the fundamental interface for all log operations.
type Logger interface {
	// Info logs an info event.
	Info(msg string, fields ...map[string]interface{})
}

// nolint: gochecknoglobals
var defaultLogger Logger = &noopLogger{}

type noopLogger struct{}

func (*noopLogger) Info(msg string, fields ...map[string]interface{}) {}
