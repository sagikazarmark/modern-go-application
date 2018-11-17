package greetingadapter

import (
	"testing"

	"github.com/go-kit/kit/log/level"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
	"github.com/stretchr/testify/assert"
)

type loggerStub struct {
	keyvals []interface{}
}

func (l *loggerStub) Log(keyvals ...interface{}) error {
	l.keyvals = keyvals

	return nil
}

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
			kitlogger := &loggerStub{}
			logger := NewLogger(kitlogger)

			test.logFunc(logger, "message: %s", name)

			expected := []interface{}{"level", test.level, "msg", "message: " + name}

			assert.Equal(t, expected, kitlogger.keyvals)
		})
	}
}

func TestLogger_WithFields(t *testing.T) {
	kitlogger := &loggerStub{}

	var logger greeting.Logger = NewLogger(kitlogger)

	logger = logger.WithFields(map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	})

	logger.Debugf("message")

	// Testing expected keys as a map to avoid problems with the unordered nature of maps
	expected := map[string]interface{}{
		"level": level.DebugValue(),
		"key1":  "value1",
		"key2":  "value2",
		"msg":   "message",
	}

	actual := make(map[string]interface{}, len(kitlogger.keyvals)/2)

	for i := 0; i < len(kitlogger.keyvals); i += 2 {
		actual[kitlogger.keyvals[i].(string)] = kitlogger.keyvals[i+1]
	}

	assert.Equal(t, expected, actual)
}
