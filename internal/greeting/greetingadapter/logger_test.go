package greetingadapter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/InVisionApp/go-logger/shims/testlog"
	"github.com/go-kit/kit/log/level"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
	"github.com/stretchr/testify/assert"
)

func TestLogger_Levels(t *testing.T) {
	tests := map[string]struct {
		logFunc func(logger *Logger, msg string, args ...interface{})
		level   level.Value
	}{
		"debug": {
			logFunc: (*Logger).Debugf,
			level:   level.DebugValue(),
		},
		"info": {
			logFunc: (*Logger).Infof,
			level:   level.InfoValue(),
		},
		"warn": {
			logFunc: (*Logger).Warnf,
			level:   level.WarnValue(),
		},
		"error": {
			logFunc: (*Logger).Errorf,
			level:   level.ErrorValue(),
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(name, func(t *testing.T) {
			testlogger := testlog.New()
			logger := NewLogger(testlogger)

			test.logFunc(logger, "message: %s", name)

			assert.Equal(t, 1, testlogger.CallCount())
			assert.Equal(t, fmt.Sprintf("[%s] %s \n", strings.ToUpper(name), "message: "+name), string(testlogger.Bytes()))
		})
	}
}

func TestLogger_WithFields(t *testing.T) {
	testlogger := testlog.New()

	var logger greeting.Logger = NewLogger(testlogger)

	logger = logger.WithFields(map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	})

	logger.Debugf("message")

	assert.Equal(t, 1, testlogger.CallCount())

	line := string(testlogger.Bytes())
	assert.Contains(t, line, "[DEBUG] message")
	assert.Contains(t, line, "key1=value1")
	assert.Contains(t, line, "key2=value2")
}
