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
	"github.com/InVisionApp/go-health/handlers"
	"github.com/cloudflare/tableflip"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/conf"
	"github.com/goph/emperror"
	"github.com/goph/emperror/errorlog"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/sagikazarmark/modern-go-application/internal"
	"github.com/sagikazarmark/modern-go-application/internal/helloworld"
	"github.com/sagikazarmark/modern-go-application/internal/helloworld/driver/web"
	"github.com/sagikazarmark/modern-go-application/internal/platform/database"
	"github.com/sagikazarmark/modern-go-application/internal/platform/invisionkitlog"
	"github.com/sagikazarmark/modern-go-application/internal/platform/jaeger"
	"github.com/sagikazarmark/modern-go-application/internal/platform/log"
	"github.com/sagikazarmark/modern-go-application/internal/platform/runner"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

func main() {
	config := NewConfig()

	config.Prepare(conf.Global)

	showVersion := conf.BoolF("version", false, "Show version information")

	conf.Parse()

	if *showVersion {
		fmt.Printf("%s version %s (%s) built on %s\n", FriendlyServiceName, Version, CommitHash, BuildDate)

		os.Exit(0)
	}

	err := config.Validate()
	if err != nil {
		fmt.Println(err)

		os.Exit(3)
	}

	// Create logger
	logger := log.NewLogger(config.Log)

	// Provide some basic context to all log lines
	logger = kitlog.With(logger, "environment", config.Environment, "service", ServiceName)

	// Configure error handler
	errorHandler := errorlog.NewHandler(logger)
	defer emperror.HandleRecover(errorHandler)

	level.Info(logger).Log(
		"version", Version, "commit_hash", CommitHash, "build_date", BuildDate,
		"msg", "starting",
	)

	instrumentRouter := http.NewServeMux()

	// Configure health checker
	healthz := health.New()
	healthz.Logger = invisionkitlog.New(logger)
	instrumentRouter.Handle("/healthz", handlers.NewJSONHandlerFunc(healthz, nil))

	// Configure Prometheus
	if config.PrometheusEnabled {
		level.Debug(logger).Log("msg", "prometheus exporter enabled")

		exporter, err := prometheus.NewExporter(prometheus.Options{
			OnError: errorHandler.Handle,
		})
		if err != nil {
			panic(errors.Wrap(err, "failed to create prometheus exporter"))
		}

		view.RegisterExporter(exporter)
		instrumentRouter.Handle("/metrics", exporter)
	}

	// Trace everything in development environment
	if config.Environment == "development" {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	// Configure Jaeger
	if config.JaegerEnabled {
		level.Debug(logger).Log("msg", "jaeger exporter enabled")

		exporter, err := jaeger.NewExporter(config.Jaeger, ServiceName, errorHandler)
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
			level.Info(logger).Log("msg", "graceful reloading")

			_ = upg.Upgrade()
		}
	}()

	// Set up instrumentation server
	{
		name := "instrumentation"
		logger := kitlog.With(logger, "server", name)
		server := &http.Server{
			Handler:  instrumentRouter,
			ErrorLog: log.NewStandardLogger(level.Error(logger)),
		}

		level.Info(logger).Log("msg", "listening on address", "address", config.InstrumentAddr)

		ln, err := upg.Fds.Listen("tcp", config.InstrumentAddr)
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
	level.Debug(logger).Log("msg", "connecting to database")
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
		err = healthz.AddCheck(&health.Config{
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

	helloWorldUseCase := &helloworld.UseCase{}
	helloWorldDriver := web.NewHelloWorldDriver(
		helloWorldUseCase,
		web.Logger(logger),
		web.ErrorHandler(errorHandler),
	)

	router := internal.NewRouter(helloWorldDriver)

	// Set up HTTP server
	{
		name := "http"
		logger := kitlog.With(logger, "server", name)
		server := &http.Server{
			Handler: &ochttp.Handler{
				Handler: router,
			},
			ErrorLog: log.NewStandardLogger(level.Error(logger)),
		}

		level.Info(logger).Log("msg", "listening on address", "address", config.HTTPAddr)

		ln, err := upg.Fds.Listen("tcp", config.HTTPAddr)
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
					level.Info(logger).Log("msg", "captured signal", "signal", sig)
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
				//level.Info(logger).Log("msg", "upgrading")
				upg.Stop()
			},
		)
	}

	if err := healthz.Start(); err != nil {
		panic(errors.Wrap(err, "failed to start health checker"))
	}

	err = group.Run()
	if err != nil {
		errorHandler.Handle(err)
	}
}
