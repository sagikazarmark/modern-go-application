package main

import (
	"context"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/conf"
	"github.com/goph/emperror"
	"github.com/goph/emperror/errorlog"
	"github.com/pkg/errors"
	"github.com/sagikazarmark/go-service-project-boilerplate/internal"
	"github.com/sagikazarmark/go-service-project-boilerplate/internal/helloworld"
	"github.com/sagikazarmark/go-service-project-boilerplate/internal/helloworld/driver/web"
	"github.com/sagikazarmark/go-service-project-boilerplate/internal/platform/database"
	"github.com/sagikazarmark/go-service-project-boilerplate/internal/platform/jaeger"
	"github.com/sagikazarmark/go-service-project-boilerplate/internal/platform/log"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

func main() {
	config := NewConfig()

	config.Prepare(conf.Global)
	conf.Parse()

	appCtx := NewAppContext(config)

	if config.ShowVersion {
		fmt.Printf(
			"%s version %s (%s) built on %s",
			appCtx.FriendlyName,
			appCtx.Build.Version,
			appCtx.Build.CommitHash,
			appCtx.Build.Date,
		)

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
	logger = kitlog.With(
		logger,
		"environment", config.Environment,
		"service", appCtx.Name,
	)

	// Configure error handler
	errorHandler := errorlog.NewHandler(logger)

	defer emperror.HandleRecover(errorHandler)

	level.Info(logger).Log(
		"version", appCtx.Build.Version,
		"commit_hash", appCtx.Build.CommitHash,
		"build_date", appCtx.Build.Date,
		"msg", "starting",
	)

	// Trace everything in development environment
	if config.Environment == "development" {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	// Configure Jaeger
	if config.JaegerEnabled {
		level.Debug(logger).Log("msg", "jaeger exporter enabled")

		exporter, err := jaeger.NewExporter(config.Jaeger, appCtx.Name, errorHandler)
		if err != nil {
			panic(err)
		}

		trace.RegisterExporter(exporter)
	}

	// Connect to the database
	level.Debug(logger).Log("msg", "connecting to database")
	db, err := database.NewConnection(config.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	helloWorld := &helloworld.UseCase{}
	helloWorldDriver := web.NewHelloWorldDriver(
		helloWorld,
		web.Logger(logger),
		web.ErrorHandler(errorHandler),
	)

	router := internal.NewRouter(helloWorldDriver)

	httpErrorLog := stdlog.New(kitlog.NewStdlibAdapter(level.Error(logger)), "", 0)

	httpServer := &http.Server{
		Addr: config.HTTPAddr,
		Handler: &ochttp.Handler{
			Handler: router,
		},
		ErrorLog: httpErrorLog,
	}
	defer httpServer.Close()

	httpServerChan := make(chan error, 1)
	go startHTTPServer(httpServer, httpServerChan, logger)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-signalChan:
		level.Info(logger).Log("msg", fmt.Sprintf("captured %v signal", sig))

	case err := <-httpServerChan:
		if err != nil && err != http.ErrServerClosed {
			errorHandler.Handle(errors.Wrap(err, "private API server crashed"))
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancel()

	go stopHTTPServer(ctx, httpServer, errorHandler)
}
