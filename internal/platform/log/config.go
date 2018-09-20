package log

import (
	"github.com/pkg/errors"
)

// Log level constants
const (
	debugLevel   = "debug"
	infoLevel    = "info"
	warnLevel    = "warn"
	warningLevel = "warning"
	errorLevel   = "error"
)

// Config holds details necessary for logging.
type Config struct {
	// Format specifies the output log format.
	// Accepted values are: json, logfmt
	Format string

	// Level is the minimum log level that should appear on the output.
	Level string
}

// NewConfig returns a new Config instance with some defaults.
func NewConfig() Config {
	return Config{
		Format: "json",
		Level:  infoLevel,
	}
}

// Validate validates the configuration.
func (c Config) Validate() error {
	if c.Format == "" {
		return errors.New("log format is required")
	}

	if c.Format != "json" && c.Format != "logfmt" {
		return errors.New("invalid log format: " + c.Format)
	}

	return nil
}
