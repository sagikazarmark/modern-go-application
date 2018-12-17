package greetingworker_test

import (
	"context"
	"testing"

	"github.com/goph/logur"
	. "github.com/sagikazarmark/modern-go-application/internal/greetingworker"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker/greetingworkeradapter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelloWorldEventLogger_SaidHello(t *testing.T) {
	logger := logur.NewTestLogger()

	eventLogger := NewHelloWorldEventLogger(greetingworkeradapter.NewLogger(logger))

	event := SaidHello{
		Message: "Hello, World!",
	}

	err := eventLogger.SaidHello(context.Background(), event)
	require.NoError(t, err)

	lastLogEvent := logger.LastEvent()
	assert.Equal(t, "said hello", lastLogEvent.Line)
	assert.Equal(t, logur.Info, lastLogEvent.Level)
	assert.Equal(t, event.Message, lastLogEvent.Fields["message"])
}
