// Package log configures a new logger for an application.
package log

import (
	"os"

	"github.com/goph/logur"
	"github.com/goph/logur/adapters/logrusadapter"
	"github.com/sirupsen/logrus"
)

// NewLogger creates a new logger.
func NewLogger(config Config) logur.Logger {
	logger := logrus.New()

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:             config.NoColor,
		EnvironmentOverrideColors: true,
	})

	switch config.Format {
	case "logfmt":
		// Already the default

	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	if level, err := logrus.ParseLevel(config.Level); err == nil {
		logger.SetLevel(level)
	}

	return logrusadapter.New(logger)
}

// WithFields returns a new contextual logger instance with context added to it.
func WithFields(logger logur.Logger, fields map[string]interface{}) logur.Logger {
	return logur.WithFields(logger, fields)
}
