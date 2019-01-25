package greetingworker

import (
	"context"
)

// GreeterEventLogger logs hello world events.
type GreeterEventLogger struct {
	logger Logger
}

// NewGreeterEventLogger returns a new GreeterEventLogger.
func NewGreeterEventLogger(logger Logger) *GreeterEventLogger {
	return &GreeterEventLogger{
		logger: logger,
	}
}

// SaidHello logs a SaidHello event.
func (e *GreeterEventLogger) SaidHelloTo(ctx context.Context, event SaidHello) error {
	e.logger.Info("said hello to someone", map[string]interface{}{"message": event.Message, "who": event.Who})

	return nil
}
