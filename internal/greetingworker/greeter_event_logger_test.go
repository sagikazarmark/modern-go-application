package greetingworker_test

import (
	"context"
	"testing"

	"github.com/goph/logur"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/sagikazarmark/modern-go-application/internal/greetingworker"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker/greetingworkeradapter"
)

func TestGreeterEventLogger_SaidHelloTo(t *testing.T) {
	logger := logur.NewTestLogger()

	eventLogger := NewGreeterEventLogger(greetingworkeradapter.NewLogger(logger))

	event := SaidHello{
		Message: "Hello, World!",
		Who:     "John",
	}

	err := eventLogger.SaidHelloTo(context.Background(), event)
	require.NoError(t, err)

	lastLogEvent := logger.LastEvent()
	assert.Equal(t, "said hello to someone", lastLogEvent.Line)
	assert.Equal(t, logur.Info, lastLogEvent.Level)
	assert.Equal(t, event.Message, lastLogEvent.Fields["message"])
	assert.Equal(t, event.Who, lastLogEvent.Fields["who"])
}
