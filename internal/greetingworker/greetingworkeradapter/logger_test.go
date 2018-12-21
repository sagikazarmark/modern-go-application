package greetingworkeradapter

import (
	"fmt"
	"testing"

	"github.com/goph/logur"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker"
	"github.com/stretchr/testify/assert"
)

func TestLogger_Levels(t *testing.T) {
	tests := map[string]struct {
		logFunc func(logger *Logger, msg string, fields map[string]interface{})
	}{
		"info": {
			logFunc: (*Logger).Info,
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(name, func(t *testing.T) {
			testLogger := logur.NewTestLogger()
			logger := NewLogger(testLogger)

			test.logFunc(logger, fmt.Sprintf("message: %s", name), nil)

			assert.Equal(t, 1, testLogger.Count())
			assert.Equal(t, name, testLogger.LastEvent().Level.String())
			assert.Equal(t, "message: "+name, testLogger.LastEvent().Line)
		})
	}
}

func TestLogger_WithFields(t *testing.T) {
	testLogger := logur.NewTestLogger()

	var logger greetingworker.Logger = NewLogger(testLogger)

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	logger = logger.WithFields(fields)

	logger.Info("message", nil)

	assert.Equal(t, 1, testLogger.Count())

	lastEvent := testLogger.LastEvent()
	assert.Equal(t, logur.Info, lastEvent.Level)
	assert.Equal(t, "message", lastEvent.Line)
	assert.Equal(t, 2, len(lastEvent.Fields))

	for key, value := range lastEvent.Fields {
		assert.Equal(t, fields[key], value)
	}
}
