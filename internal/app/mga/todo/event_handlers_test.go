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

func TestLogEventHandler_MarkedAsDone(t *testing.T) {
	logger := &logur.TestLoggerFacade{}

	eventHandler := NewLogEventHandler(commonadapter.NewLogger(logger))

	event := MarkedAsDone{
		ID: "1234",
	}

	err := eventHandler.MarkedAsDone(context.Background(), event)
	require.NoError(t, err)

	logEvent := logur.LogEvent{
		Level: logur.Info,
		Line:  "todo marked as done",
		Fields: map[string]interface{}{
			"event":   "MarkedAsDone",
			"todo_id": "1234",
		},
	}

	logtesting.AssertLogEventsEqual(t, logEvent, *(logger.LastEvent()))
}
