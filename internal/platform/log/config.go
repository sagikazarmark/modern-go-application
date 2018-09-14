package log

import (
	"github.com/pkg/errors"
)

// Config holds details necessary for logging.
type Config struct {
	// Defines the log format.
	// Valid values are: json, logfmt
	Format string

	Environment string
	Debug       bool

	ServiceName string
}

// NewConfig returns a new Config instance with some defaults.
func NewConfig() Config {
	return Config{
		Format: "json",
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
