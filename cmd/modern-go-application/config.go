package main

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/sagikazarmark/modern-go-application/internal/platform/database"
	"github.com/sagikazarmark/modern-go-application/internal/platform/log"
	"github.com/sagikazarmark/modern-go-application/internal/platform/opencensus"
)

// configuration holds any kind of configuration that comes from the outside world and
// is necessary for running the application.
type configuration struct {
	// Log configuration
	Log log.Config

	// Telemetry configuration
	Telemetry struct {
		// Telemetry HTTP server address
		Addr string
	}

	// OpenCensus configuration
	Opencensus struct {
		Exporter struct {
			Enabled bool

			opencensus.ExporterConfig `mapstructure:",squash"`
		}

		Trace opencensus.TraceConfig
	}

	// App configuration
	App appConfig

	// Database connection information
	Database database.Config
}

// Process post-processes configuration after loading it.
func (configuration) Process() error {
	return nil
}

// Validate validates the configuration.
func (c configuration) Validate() error {
	if c.Telemetry.Addr == "" {
		return errors.New("telemetry http server address is required")
	}

	if err := c.App.Validate(); err != nil {
		return err
	}

	if err := c.Database.Validate(); err != nil {
		return err
	}

	return nil
}

// appConfig represents the application related configuration.
type appConfig struct {
	// HTTP server address
	// nolint: golint, stylecheck
	HttpAddr string

	// GRPC server address
	GrpcAddr string

	// Storage is the storage backend of the application
	Storage string
}

// Validate validates the configuration.
func (c appConfig) Validate() error {
	if c.HttpAddr == "" {
		return errors.New("http app server address is required")
	}

	if c.GrpcAddr == "" {
		return errors.New("grpc app server address is required")
	}

	if c.Storage != "inmemory" && c.Storage != "database" {
		return errors.New("app storage must be inmemory or database")
	}

	return nil
}

// configure configures some defaults in the Viper instance.
func configure(v *viper.Viper, p *pflag.FlagSet) {
	// Viper settings
	v.AddConfigPath(".")
	v.AddConfigPath("$CONFIG_DIR/")

	// Environment variable settings
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	// Global configuration
	v.SetDefault("shutdownTimeout", 15*time.Second)
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		v.SetDefault("no_color", true)
	}

	// Log configuration
	v.SetDefault("log.format", "json")
	v.SetDefault("log.level", "info")
	v.RegisterAlias("log.noColor", "no_color")

	// Telemetry configuration
	p.String("telemetry-addr", ":10000", "Telemetry HTTP server address")
	_ = v.BindPFlag("telemetry.addr", p.Lookup("telemetry-addr"))
	v.SetDefault("telemetry.addr", ":10000")

	// OpenCensus configuration
	v.SetDefault("opencensus.exporter.enabled", false)
	_ = v.BindEnv("opencensus.exporter.address")
	_ = v.BindEnv("opencensus.exporter.insecure")
	_ = v.BindEnv("opencensus.exporter.reconnectPeriod")
	v.SetDefault("opencensus.trace.sampling.sampler", "never")
	v.SetDefault("opencensus.prometheus.enabled", false)

	// App configuration
	p.String("http-addr", ":8000", "App HTTP server address")
	_ = v.BindPFlag("app.httpAddr", p.Lookup("http-addr"))
	v.SetDefault("app.httpAddr", ":8000")

	p.String("grpc-addr", ":8001", "App GRPC server address")
	_ = v.BindPFlag("app.grpcAddr", p.Lookup("grpc-addr"))
	v.SetDefault("app.grpcAddr", ":8001")

	v.SetDefault("app.storage", "inmemory")

	// Database configuration
	_ = v.BindEnv("database.host")
	v.SetDefault("database.port", 3306)
	_ = v.BindEnv("database.user")
	_ = v.BindEnv("database.pass")
	_ = v.BindEnv("database.name")
	v.SetDefault("database.params", map[string]string{
		"collation": "utf8mb4_general_ci",
	})
}
