package todo_test

import (
	"context"
	"testing"

	"github.com/goph/logur"
	"github.com/goph/logur/logtesting"
	"github.com/stretchr/testify/require"

	. "github.com/sagikazarmark/modern-go-application/internal/todo"
	"github.com/sagikazarmark/modern-go-application/internal/todo/todoadapter"
)

func TestLogEventHandler_MarkedAsDone(t *testing.T) {
	logger := logur.NewTestLogger()

	eventHandler := NewLogEventHandler(todoadapter.NewLogger(logger))

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
