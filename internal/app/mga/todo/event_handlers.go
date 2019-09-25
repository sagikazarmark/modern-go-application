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

// MarkedAsDone logs a MarkedAsDone event.
func (h LogEventHandler) MarkedAsDone(ctx context.Context, event MarkedAsDone) error {
	logger := h.logger.WithContext(ctx)

	logger.Info("todo marked as done", map[string]interface{}{
		"event":   "MarkedAsDone",
		"todo_id": event.ID,
	})

	return nil
}
