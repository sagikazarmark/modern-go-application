package todo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"logur.dev/logur"
	"logur.dev/logur/logtesting"

	. "github.com/sagikazarmark/modern-go-application/internal/app/mga/todo"
	"github.com/sagikazarmark/modern-go-application/internal/common/commonadapter"
)

func TestLogEventHandler_MarkedAsComplete(t *testing.T) {
	logger := &logur.TestLoggerFacade{}

	eventHandler := NewLogEventHandler(commonadapter.NewLogger(logger))

	event := MarkedAsComplete{
		ID: "1234",
	}

	err := eventHandler.MarkedAsComplete(context.Background(), event)
	require.NoError(t, err)

	logEvent := logur.LogEvent{
		Level: logur.Info,
		Line:  "todo marked as complete",
		Fields: map[string]interface{}{
			"event":   "MarkedAsComplete",
			"todo_id": "1234",
		},
	}

	logtesting.AssertLogEventsEqual(t, logEvent, *(logger.LastEvent()))
}
