package main

import (
	"errors"
	"time"

	"github.com/goph/conf"
	"github.com/sagikazarmark/go-service-project-boilerplate/internal/platform/database"
	"github.com/sagikazarmark/go-service-project-boilerplate/internal/platform/log"
)

// Config holds any kind of configuration that comes from the outside world and
// is necessary for running the application.
type Config struct {
	// Show the application build information.
	ShowVersion bool

	// Meaningful values are recommended (eg. production, development, staging, release/123, etc)
	//
	// "development" is treated special: address types will use the loopback interface as default when none is defined.
	// This is useful when developing locally and listening on all interfaces requires elevated rights.
	Environment string

	// Turns on some debug functionality: more verbose logs, exposed pprof, expvar and net trace in the debug server.
	Debug bool

	// Timeout for graceful shutdown
	ShutdownTimeout time.Duration

	// Log configuration
	Log log.Config

	// HTTP address
	HTTPAddr string

	// Database connection information
	Database database.Config
}

// NewConfig returns a Config struct with sane defaults.
func NewConfig() Config {
	return Config{
		Environment:     "production",
		ShutdownTimeout: 15 * time.Second,
		Log:             log.NewConfig(),
		HTTPAddr:        ":8000",
		Database:        database.NewConfig(),
	}
}

// Validate validates the configuration.
func (c Config) Validate() error {
	if c.Environment == "" {
		return errors.New("environment is required")
	}

	if c.HTTPAddr == "" {
		return errors.New("http address is required")
	}

	if err := c.Log.Validate(); err != nil {
		return err
	}

	if err := c.Database.Validate(); err != nil {
		return err
	}

	return nil
}

// Prepare prepares the configuration to be populated from various sources
// (determined by the console nature of the application).
func (c *Config) Prepare(conf *conf.Configurator) {
	conf.BoolVarF(&c.ShowVersion, "version", false, "Show version information")

	// General configuration
	conf.StringVar(&c.Environment, "environment", c.Environment, "Application environment")
	conf.BoolVar(&c.Debug, "debug", c.Debug, "Turns on debug functionality")
	conf.DurationVarF(&c.ShutdownTimeout, "shutdown-timeout", c.ShutdownTimeout, "Timeout for graceful shutdown")

	// Log configuration
	conf.StringVar(&c.Log.Format, "log-format", c.Log.Format, "Defines the log format (json or logfmt)")

	conf.StringVarF(&c.HTTPAddr, "http-addr", c.HTTPAddr, "HTTP address")

	// Database configuration
	conf.StringVar(&c.Database.Host, "db-host", c.Database.Host, "Database host")
	conf.IntVar(&c.Database.Port, "db-port", c.Database.Port, "Database port")
	conf.StringVar(&c.Database.User, "db-user", c.Database.User, "Database user")
	conf.StringVar(&c.Database.Pass, "db-pass", c.Database.Pass, "Database password")
	conf.StringVar(&c.Database.Name, "db-name", c.Database.Name, "Database name")
	conf.QueryStringVar(&c.Database.Params, "db-params", c.Database.Params, "Database params")
}

// ApplyContext updates configuration values based on the application context.
func (c *Config) ApplyContext(ctx Context) {
	c.Log.Debug = ctx.Debug
	c.Log.Environment = ctx.Environment
	c.Log.ServiceName = ctx.Name
}
