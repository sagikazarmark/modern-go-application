package greetingworker

import (
	"context"
)

// SayHelloEventLogger logs hello world events.
type SayHelloEventLogger struct {
	logger Logger
}

// NewSayHelloEventLogger returns a new SayHelloEventLogger.
func NewSayHelloEventLogger(logger Logger) *SayHelloEventLogger {
	return &SayHelloEventLogger{
		logger: logger,
	}
}

// SaidHelloTo logs a SaidHelloTo event.
func (e *SayHelloEventLogger) SaidHelloTo(ctx context.Context, event SaidHelloTo) error {
	e.logger.Info("said hello to someone", map[string]interface{}{"message": event.Message, "who": event.Who})

	return nil
}
