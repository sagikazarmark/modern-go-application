package greetingadapter

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

// Logger wraps a go-kit logger and exposes it under a custom interface.
type Logger struct {
	logger log.Logger
}

// NewLogger returns a new Logger instance.
func NewLogger(logger log.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	_ = level.Debug(l.logger).Log("msg", fmt.Sprintf(msg, args...))
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	_ = level.Info(l.logger).Log("msg", fmt.Sprintf(msg, args...))
}

func (l *Logger) Warnf(msg string, args ...interface{}) {
	_ = level.Warn(l.logger).Log("msg", fmt.Sprintf(msg, args...))
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	_ = level.Error(l.logger).Log("msg", fmt.Sprintf(msg, args...))
}

func (l *Logger) WithFields(fields map[string]interface{}) greeting.Logger {
	keyvals := make([]interface{}, len(fields)*2)
	i := 0

	for key, value := range fields {
		keyvals[i] = key
		keyvals[i+1] = value

		i += 2
	}

	return &Logger{
		logger: log.With(l.logger, keyvals...),
	}
}
