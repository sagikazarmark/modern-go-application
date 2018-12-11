package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/InVisionApp/go-health"
	"github.com/InVisionApp/go-health/checkers"
	"github.com/cloudflare/tableflip"
	"github.com/goph/emperror"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/sagikazarmark/modern-go-application/internal"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingadapter"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingdriver"
	"github.com/sagikazarmark/modern-go-application/internal/platform/buildinfo"
	"github.com/sagikazarmark/modern-go-application/internal/platform/database"
	"github.com/sagikazarmark/modern-go-application/internal/platform/errorhandler"
	"github.com/sagikazarmark/modern-go-application/internal/platform/healthcheck"
	"github.com/sagikazarmark/modern-go-application/internal/platform/jaeger"
	"github.com/sagikazarmark/modern-go-application/internal/platform/log"
	"github.com/sagikazarmark/modern-go-application/internal/platform/prometheus"
	"github.com/sagikazarmark/modern-go-application/internal/platform/runner"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// nolint: gochecknoinits
func init() {
	pflag.Bool("version", false, "Show version information")
	pflag.Bool("dump-config", false, "Dump configuration to the console")
}

func main() {
	Configure(viper.GetViper(), pflag.CommandLine)

	pflag.Parse()

	if viper.GetBool("version") {
		fmt.Printf("%s version %s (%s) built on %s\n", FriendlyServiceName, version, commitHash, buildDate)

		os.Exit(0)
	}

	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); err != nil && !ok {
		panic(errors.Wrap(err, "failed to read configuration"))
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(errors.Wrap(err, "failed to unmarshal configuration"))
	}

	err = config.Validate()
	if err != nil {
		fmt.Println(err)

		os.Exit(3)
	}

	if viper.GetBool("dump-config") {
		fmt.Printf("%+v\n", config)

		os.Exit(0)
	}

	// Create logger
	logger := log.NewLogger(config.Log)

	// Provide some basic context to all log lines
	logger = logger.WithFields(log.Fields{"environment": config.Environment, "service": ServiceName})

	// Configure error handler
	errorHandler := errorhandler.New(logger)
	defer emperror.HandleRecover(errorHandler)

	buildInfo := buildinfo.New(version, commitHash, buildDate)

	logger.WithFields(log.Fields(buildInfo.Fields())).Info("starting application")

	instrumentationRouter := http.NewServeMux()
	instrumentationRouter.Handle("/version", buildinfo.Handler(buildInfo))

	// Configure health checker
	healthChecker := healthcheck.New(logger)
	instrumentationRouter.Handle("/healthz", healthcheck.NewHTTPHandler(healthChecker))

	// Configure Prometheus
	if config.Instrumentation.Prometheus.Enabled {
		logger.Info("prometheus exporter enabled")

		exporter, err := prometheus.NewExporter(config.Instrumentation.Prometheus.Config, errorHandler)
		if err != nil {
			panic(err)
		}

		view.RegisterExporter(exporter)
		instrumentationRouter.Handle("/metrics", exporter)
	}

	// Trace everything in development environment
	if config.Environment == "development" {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	// Configure Jaeger
	if config.Instrumentation.Jaeger.Enabled {
		logger.Info("jaeger exporter enabled")

		exporter, err := jaeger.NewExporter(config.Instrumentation.Jaeger.Config, ServiceName, errorHandler)
		if err != nil {
			panic(err)
		}

		trace.RegisterExporter(exporter)
	}

	// Configure graceful restart
	upg, _ := tableflip.New(tableflip.Options{})

	var group run.Group

	// Do an upgrade on SIGHUP
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGHUP)
		for range ch {
			logger.Info("graceful reloading")

			_ = upg.Upgrade()
		}
	}()

	// Set up instrumentation server
	{
		name := "instrumentation"
		logger := logger.WithFields(log.Fields{"server": name})
		server := &http.Server{
			Handler:  instrumentationRouter,
			ErrorLog: log.NewStandardErrorLogger(logger),
		}

		logger.WithFields(log.Fields{"address": config.Instrumentation.Addr}).Info("listening on address")

		ln, err := upg.Fds.Listen("tcp", config.Instrumentation.Addr)
		if err != nil {
			panic(err)
		}

		r := &runner.Server{
			Server:          server,
			Listener:        ln,
			ShutdownTimeout: config.ShutdownTimeout,
			Logger:          logger,
			ErrorHandler:    emperror.HandlerWith(errorHandler, "server", name),
		}

		group.Add(r.Start, r.Stop)
	}

	// Connect to the database
	logger.Info("connecting to database")
	db, err := database.NewConnection(config.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Register database health check
	{
		check, err := checkers.NewSQL(&checkers.SQLConfig{Pinger: db})
		if err != nil {
			panic(errors.Wrap(err, "failed to create db health checker"))
		}
		err = healthChecker.AddCheck(&health.Config{
			Name:     "database",
			Checker:  check,
			Interval: time.Duration(3) * time.Second,
			Fatal:    true,
		})
		if err != nil {
			panic(errors.Wrap(err, "failed to add health checker"))
		}
	}

	// Register HTTP stat views
	if err := view.Register(ochttp.DefaultServerViews...); err != nil {
		panic(errors.Wrap(err, "failed to register HTTP server stat views"))
	}

	helloWorld := greeting.NewHelloWorld(greetingadapter.NewLogger(logger))
	sayHello := greeting.NewSayHello(greetingadapter.NewLogger(logger))
	helloWorldController := greetingdriver.NewGreetingController(helloWorld, sayHello, errorHandler)

	router := internal.NewRouter(helloWorldController)

	// Set up app server
	{
		name := "app"
		logger := logger.WithFields(log.Fields{"server": name})
		server := &http.Server{
			Handler: &ochttp.Handler{
				Handler: router,
			},
			ErrorLog: log.NewStandardErrorLogger(logger),
		}

		logger.WithFields(log.Fields{"address": config.App.Addr}).Info("listening on address")

		ln, err := upg.Fds.Listen("tcp", config.App.Addr)
		if err != nil {
			panic(err)
		}

		r := &runner.Server{
			Server:          server,
			Listener:        ln,
			ShutdownTimeout: config.ShutdownTimeout,
			Logger:          logger,
			ErrorHandler:    emperror.HandlerWith(errorHandler, "server", name),
		}

		group.Add(r.Start, r.Stop)
	}

	// Setup exit signal
	{
		ch := make(chan os.Signal, 1)

		group.Add(
			func() error {
				signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

				sig := <-ch
				if sig != nil {
					logger.WithFields(log.Fields{"signal": sig}).Info("captured signal")
				}

				return nil
			},
			func(e error) {
				signal.Stop(ch)
				close(ch)
			},
		)
	}

	{
		group.Add(
			func() error {
				// Tell the parent we are ready
				_ = upg.Ready()

				// Wait for children to be ready
				// (or application shutdown)
				<-upg.Exit()

				return nil
			},
			func(e error) {
				upg.Stop()
			},
		)
	}

	if err := healthChecker.Start(); err != nil {
		panic(errors.Wrap(err, "failed to start health checker"))
	}

	err = group.Run()
	if err != nil {
		errorHandler.Handle(err)
	}
}
