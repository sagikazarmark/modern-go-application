// Package log configures a new logger for an application.
package log

import (
	"os"

	"github.com/InVisionApp/go-logger"
	logrusShim "github.com/InVisionApp/go-logger/shims/logrus"
	"github.com/sirupsen/logrus"
)

// Fields is an alias to log.Fields for easier usage.
type Fields = log.Fields

// NewLogger creates a new logger.
func NewLogger(config Config) log.Logger {
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
		logrus.SetLevel(level)
	}

	return logrusShim.New(logger)
}
