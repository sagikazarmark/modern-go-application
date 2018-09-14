package app

import (
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

// LogConfig holds details necessary for logging.
type LogConfig struct {
	// Defines the log format.
	// Valid values are: json, logfmt
	Format string
}

// NewLogConfig returns a new LogConfig instance with some defaults.
func NewLogConfig() LogConfig {
	return LogConfig{
		Format: "json",
	}
}

// Validate validates the configuration.
func (c LogConfig) Validate() error {
	if c.Format == "" {
		return errors.New("log format is required")
	}

	if c.Format != "json" && c.Format != "logfmt" {
		return errors.New("invalid log format: " + c.Format)
	}

	return nil
}

// NewLogger creates a new logger.
func NewLogger(config LogConfig, appCtx Context) (log.Logger, error) {
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

	// Provide some basic context to all log lines
	logger = log.With(
		logger,
		"environment", appCtx.Environment,
		"service", appCtx.Name,
	)

	// Fallback to Info level
	logger = level.NewInjector(logger, level.InfoValue())

	// Only log debug level messages if debug mode is turned on
	if !appCtx.Debug {
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	return logger, nil
}
