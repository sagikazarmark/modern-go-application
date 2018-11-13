// Package logger configures a new logger for an application.
package log

import (
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// NewLogger creates a new logger.
func NewLogger(config Config) log.Logger {
	var logger log.Logger

	w := log.NewSyncWriter(os.Stdout)

	switch config.Format {
	case "logfmt":
		logger = log.NewLogfmtLogger(w)

	case "json":
		logger = log.NewJSONLogger(w)

	default: // Fall back to JSON logger
		logger = log.NewJSONLogger(w)

		level.Warn(logger).Log("msg", "unsupported log format", "format", config.Format)
	}

	// Fallback to Info level
	logger = level.NewInjector(logger, level.InfoValue())

	// Set log level
	var levelOption level.Option
	switch strings.ToLower(config.Level) {
	case debugLevel:
		levelOption = level.AllowDebug()

	case infoLevel, "": // Info is the default level
		levelOption = level.AllowInfo()

	case warnLevel, warningLevel:
		levelOption = level.AllowWarn()

	case errorLevel:
		levelOption = level.AllowError()

	default: // Info is the default level
		levelOption = level.AllowInfo()
	}

	logger = level.NewFilter(logger, levelOption)

	return logger
}
