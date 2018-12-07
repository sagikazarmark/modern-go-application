package log

// Level type
type Level uint32

// These are the different logging levels.
const (
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel Level = iota

	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel

	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel

	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
)
