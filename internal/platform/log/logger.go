// Package logger configures a new logger for an application.
package log

import (
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// NewLogger creates a new logger.
func NewLogger(config Config) (log.Logger, error) {
	var logger log.Logger

	w := log.NewSyncWriter(os.Stdout)

	switch config.Format {
	case "logfmt":
		logger = log.NewLogfmtLogger(w)

	case "json":
		logger = log.NewJSONLogger(w)

	default:
		return nil, fmt.Errorf("unsupported log format: %s", config.Format)
	}

	// Fallback to Info level
	logger = level.NewInjector(logger, level.InfoValue())

	// Only log debug level messages if debug mode is turned on
	if !config.Debug {
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	return logger, nil
}
