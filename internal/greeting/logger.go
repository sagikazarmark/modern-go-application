package greeting

// Logger is the fundamental interface for all log operations.
type Logger interface {
	// Debugf logs a debug event and optionally formats the message.
	Debugf(msg string, args ...interface{})

	// Infof logs an info event and optionally formats the message.
	Infof(msg string, args ...interface{})

	// Warnf logs a warning event and optionally formats the message.
	Warnf(msg string, args ...interface{})

	// Errorf logs an error event and optionally formats the message.
	Errorf(msg string, args ...interface{})

	// WithFields annotates a logger with some context.
	WithFields(fields map[string]interface{}) Logger
}

type nopLogger struct{}

func (*nopLogger) Debugf(msg string, args ...interface{}) {}

func (*nopLogger) Infof(msg string, args ...interface{}) {}

func (*nopLogger) Warnf(msg string, args ...interface{}) {}

func (*nopLogger) Errorf(msg string, args ...interface{}) {}

func (l *nopLogger) WithFields(fields map[string]interface{}) Logger {
	return l
}

// NewNopLogger returns a logger that doesn't do anything.
func NewNopLogger() Logger {
	return &nopLogger{}
}
