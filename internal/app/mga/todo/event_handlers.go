package todo

import (
	"context"
)

// LogEventHandler handles todo events and logs them.
type LogEventHandler struct {
	logger Logger
}

// NewLogEventHandler returns a new LogEventHandler instance.
func NewLogEventHandler(logger Logger) LogEventHandler {
	return LogEventHandler{
		logger: logger,
	}
}

// MarkedAsComplete logs a MarkedAsComplete event.
func (h LogEventHandler) MarkedAsComplete(ctx context.Context, event MarkedAsComplete) error {
	logger := h.logger.WithContext(ctx)

	logger.Info("todo marked as complete", map[string]interface{}{
		"event":   "MarkedAsComplete",
		"todo_id": event.ID,
	})

	return nil
}
