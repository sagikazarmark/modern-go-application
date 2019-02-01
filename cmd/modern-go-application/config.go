package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/sagikazarmark/modern-go-application/internal/platform/database"
	"github.com/sagikazarmark/modern-go-application/internal/platform/jaeger"
	"github.com/sagikazarmark/modern-go-application/internal/platform/log"
	"github.com/sagikazarmark/modern-go-application/internal/platform/prometheus"
	"github.com/sagikazarmark/modern-go-application/internal/platform/redis"
	"github.com/sagikazarmark/modern-go-application/internal/platform/watermill"
)

// Config holds any kind of configuration that comes from the outside world and
// is necessary for running the application.
type Config struct {
	// Meaningful values are recommended (eg. production, development, staging, release/123, etc)
	Environment string

	// Turns on some debug functionality
	Debug bool

	// Timeout for graceful shutdown
	ShutdownTimeout time.Duration

	// Log configuration
	Log log.Config

	// Instrumentation configuration
	Instrumentation InstrumentationConfig

	// App configuration
	App struct {
		// App server address
		Addr string
	}

	// Database connection information
	Database database.Config

	// Redis configuration
	Redis redis.Config

	// Watermill configuration
	Watermill struct {
		RouterConfig watermill.RouterConfig
	}
}

// Validate validates the configuration.
func (c Config) Validate() error {
	if c.Environment == "" {
		return errors.New("environment is required")
	}

	if err := c.Instrumentation.Validate(); err != nil {
		return err
	}

	if c.App.Addr == "" {
		return errors.New("app server address is required")
	}

	if err := c.Database.Validate(); err != nil {
		return err
	}

	// Uncomment to enable redis config validation
	// if err := c.Redis.Validate(); err != nil {
	// 	return err
	// }

	return nil
}

// InstrumentationConfig represents the instrumentation related configuration.
type InstrumentationConfig struct {
	// Instrumentation HTTP server address
	Addr string

	// Prometheus configuration
	Prometheus struct {
		Enabled           bool
		prometheus.Config `mapstructure:",squash"`
	}

	// Jaeger configuration
	Jaeger struct {
		Enabled       bool
		jaeger.Config `mapstructure:",squash"`
	}
}

// Validate validates the configuration.
func (c InstrumentationConfig) Validate() error {
	if c.Addr == "" {
		return errors.New("instrumentation http server address is required")
	}

	if c.Jaeger.Enabled {
		if err := c.Jaeger.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Configure configures some defaults in the Viper instance.
func Configure(v *viper.Viper, p *pflag.FlagSet) {
	v.AllowEmptyEnv(true)
	v.AddConfigPath(".")
	v.AddConfigPath(fmt.Sprintf("$%s_CONFIG_DIR/", strings.ToUpper(EnvPrefix)))
	p.Init(FriendlyServiceName, pflag.ExitOnError)
	pflag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", FriendlyServiceName)
		pflag.PrintDefaults()
	}
	_ = v.BindPFlags(p)

	v.SetEnvPrefix(EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	// Application constants
	v.Set("serviceName", ServiceName)

	// Global configuration
	v.SetDefault("environment", "production")
	v.SetDefault("debug", false)
	v.SetDefault("shutdownTimeout", 15*time.Second)
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		v.SetDefault("no_color", true)
	}

	// Log configuration
	v.SetDefault("log.format", "json")
	v.SetDefault("log.level", "info")
	v.RegisterAlias("log.noColor", "no_color")

	// Instrumentation configuration
	p.String("instrumentation.addr", ":10000", "Instrumentation HTTP server address")
	v.SetDefault("instrumentation.addr", ":10000")

	v.SetDefault("instrumentation.prometheus.enabled", false)
	v.SetDefault("instrumentation.jaeger.enabled", false)
	_ = v.BindEnv("instrumentation.jaeger.collectorEndpoint")
	v.SetDefault("instrumentation.jaeger.agentEndpoint", "localhost:6831")
	v.RegisterAlias("instrumentation.jaeger.serviceName", "serviceName")
	_ = v.BindEnv("instrumentation.jaeger.username")
	_ = v.BindEnv("instrumentation.jaeger.password")

	// App configuration
	p.String("app.addr", ":8000", "App HTTP server address")
	v.SetDefault("app.addr", ":8000")

	// Database configuration
	_ = v.BindEnv("database.host")
	v.SetDefault("database.port", 3306)
	_ = v.BindEnv("database.user")
	_ = v.BindEnv("database.pass")
	_ = v.BindEnv("database.name")
	v.SetDefault("database.params", map[string]string{
		"charset": "utf8mb4",
	})

	// Redis configuration
	_ = v.BindEnv("redis.host")
	v.SetDefault("redis.port", 6379)
	_ = v.BindEnv("redis.password")

	// Watermill configuration
	v.RegisterAlias("watermill.routerConfig.closeTimeout", "shutdownTimeout")
}
