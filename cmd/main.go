package main

import (
	"context"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"sync"
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
	"github.com/pkg/errors"
	"github.com/sagikazarmark/modern-go-application/internal"
	"github.com/sagikazarmark/modern-go-application/internal/helloworld"
	"github.com/sagikazarmark/modern-go-application/internal/helloworld/driver/web"
	"github.com/sagikazarmark/modern-go-application/internal/platform/database"
	"github.com/sagikazarmark/modern-go-application/internal/platform/invisionkitlog"
	"github.com/sagikazarmark/modern-go-application/internal/platform/jaeger"
	"github.com/sagikazarmark/modern-go-application/internal/platform/log"
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
	logger, err := log.NewLogger(config.Log)
	if err != nil {
		panic(err)
	}

	// Provide some basic context to all log lines
	logger = kitlog.With(logger, "environment", config.Environment, "service", ServiceName, "pid", os.Getpid())

	// Configure error handler
	errorHandler := errorlog.NewHandler(logger)
	defer emperror.HandleRecover(errorHandler)

	level.Info(logger).Log(
		"version", Version, "commit_hash", CommitHash, "build_date", BuildDate,
		"msg", "starting",
	)

	// Configure health checker
	healthz := health.New()
	healthz.Logger = invisionkitlog.New(logger)

	// Connect to the database
	level.Debug(logger).Log("msg", "connecting to database")
	db, err := database.NewConnection(config.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	dbCheck, err := checkers.NewSQL(&checkers.SQLConfig{Pinger: db})
	if err != nil {
		panic(errors.Wrap(err, "failed to create db health checker"))
	}
	err = healthz.AddCheck(&health.Config{
		Name:     "database",
		Checker:  dbCheck,
		Interval: time.Duration(3) * time.Second,
		Fatal:    true,
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to add health checker"))
	}

	instrumentRouter := http.NewServeMux()

	if err := healthz.Start(); err != nil {
		panic(errors.Wrap(err, "failed to start health checker"))
	}

	instrumentRouter.Handle("/healthz", handlers.NewJSONHandlerFunc(healthz, nil))

	// Configure prometheus
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

	// Graceful restart
	upg, _ := tableflip.New(tableflip.Options{})
	defer upg.Stop()

	// Do an upgrade on SIGHUP
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP)
		for range sig {
			level.Info(logger).Log("msg", "graceful reloading")

			_ = upg.Upgrade()
		}
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	// Set up instrumentation server
	instrumentLogger := kitlog.With(logger, "server", "instrumentation")
	instrumentServer := &http.Server{
		Handler:  instrumentRouter,
		ErrorLog: stdlog.New(kitlog.NewStdlibAdapter(level.Error(instrumentLogger)), "", 0),
	}
	defer instrumentServer.Close()

	instrumentServerChan := make(chan error, 1)
	go func() {
		level.Info(instrumentLogger).Log("msg", "starting server", "addr", config.InstrumentAddr)

		ln, err := upg.Fds.Listen("tcp", config.InstrumentAddr)
		if err != nil {
			panic(err)
		}

		wg.Done()

		instrumentServerChan <- instrumentServer.Serve(ln)
	}()

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

	httpLogger := kitlog.With(logger, "server", "http")
	httpServer := &http.Server{
		Handler: &ochttp.Handler{
			Handler: router,
		},
		ErrorLog: stdlog.New(kitlog.NewStdlibAdapter(level.Error(httpLogger)), "", 0),
	}
	defer httpServer.Close()

	httpServerChan := make(chan error, 1)
	go func() {
		level.Info(httpLogger).Log("msg", "starting server", "addr", config.HTTPAddr)

		ln, err := upg.Fds.Listen("tcp", config.HTTPAddr)
		if err != nil {
			panic(err)
		}

		wg.Done()

		httpServerChan <- httpServer.Serve(ln)
	}()

	wg.Wait()

	// Tell the parent we are ready
	_ = upg.Ready()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	// Wait to be replaced with a new process
	case <-upg.Exit():
		level.Info(logger).Log("msg", "upgrading")

	case sig := <-signalChan:
		level.Info(logger).Log("msg", "captured signal", "signal", sig)

	case err := <-instrumentServerChan:
		if err != nil && err != http.ErrServerClosed {
			errorHandler.Handle(emperror.With(errors.Wrap(err, "http server crashed"), "server", "instrumentation"))
		}

	case err := <-httpServerChan:
		if err != nil && err != http.ErrServerClosed {
			errorHandler.Handle(emperror.With(errors.Wrap(err, "http server crashed"), "server", "http"))
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancel()

	wg.Add(2)

	// Shut down instrumentation server
	go func() {
		level.Info(instrumentLogger).Log("msg", "shutting server down")

		err := instrumentServer.Shutdown(ctx)
		if err != nil {
			errorHandler.Handle(err)
		}

		wg.Done()
	}()

	// Shut down HTTP server
	go func() {
		level.Info(httpLogger).Log("msg", "shutting server down")

		err := httpServer.Shutdown(ctx)
		if err != nil {
			errorHandler.Handle(err)
		}

		wg.Done()
	}()

	wg.Wait()
}
