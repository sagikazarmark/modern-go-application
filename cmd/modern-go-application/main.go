package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/InVisionApp/go-health"
	"github.com/InVisionApp/go-health/checkers"
	"github.com/cloudflare/tableflip"
	"github.com/goph/emperror"
	"github.com/goph/watermillx"
	"github.com/oklog/run"
	"github.com/opencensus-integrations/ocsql"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"

	"github.com/sagikazarmark/modern-go-application/internal"
	"github.com/sagikazarmark/modern-go-application/internal/platform/buildinfo"
	"github.com/sagikazarmark/modern-go-application/internal/platform/database"
	"github.com/sagikazarmark/modern-go-application/internal/platform/errorhandler"
	"github.com/sagikazarmark/modern-go-application/internal/platform/healthcheck"
	"github.com/sagikazarmark/modern-go-application/internal/platform/jaeger"
	"github.com/sagikazarmark/modern-go-application/internal/platform/log"
	"github.com/sagikazarmark/modern-go-application/internal/platform/prometheus"
	apptrace "github.com/sagikazarmark/modern-go-application/internal/platform/trace"
	"github.com/sagikazarmark/modern-go-application/internal/platform/watermill"
)

// nolint: gochecknoinits
func init() {
	pflag.Bool("version", false, "Show version information")
	pflag.Bool("dump-config", false, "Dump configuration to the console (and exit)")
}

func main() {
	Configure(viper.GetViper(), pflag.CommandLine)

	pflag.Parse()

	if viper.GetBool("version") {
		fmt.Printf("%s version %s (%s) built on %s\n", FriendlyServiceName, version, commitHash, buildDate)

		os.Exit(0)
	}

	err := viper.ReadInConfig()
	_, configFileNotFound := err.(viper.ConfigFileNotFoundError)
	if !configFileNotFound {
		emperror.Panic(errors.Wrap(err, "failed to read configuration"))
	}

	var config Config
	err = viper.Unmarshal(&config)
	emperror.Panic(errors.Wrap(err, "failed to unmarshal configuration"))

	// Create logger (first thing after configuration loading)
	logger := log.NewLogger(config.Log)

	// Provide some basic context to all log lines
	logger = log.WithFields(logger, map[string]interface{}{"environment": config.Environment, "service": ServiceName})

	log.SetStandardLogger(logger)

	if configFileNotFound {
		logger.Warn("configuration file not found", nil)
	}

	err = config.Validate()
	if err != nil {
		logger.Error(err.Error(), nil)

		os.Exit(3)
	}

	if viper.GetBool("dump-config") {
		fmt.Printf("%+v\n", config)

		os.Exit(0)
	}

	// Configure error handler
	errorHandler := errorhandler.New(logger)
	defer emperror.HandleRecover(errorHandler)

	buildInfo := buildinfo.New(version, commitHash, buildDate)

	logger.Info("starting application", buildInfo.Fields())

	instrumentationRouter := http.NewServeMux()
	instrumentationRouter.Handle("/version", buildinfo.Handler(buildInfo))

	// Configure health checker
	healthChecker := healthcheck.New(logger)
	instrumentationRouter.Handle("/healthz", healthcheck.Handler(healthChecker))

	// Configure Prometheus
	if config.Instrumentation.Prometheus.Enabled {
		logger.Info("prometheus exporter enabled", nil)

		exporter, err := prometheus.NewExporter(config.Instrumentation.Prometheus.Config, errorHandler)
		emperror.Panic(err)

		view.RegisterExporter(exporter)
		instrumentationRouter.Handle("/metrics", exporter)
	}

	// Trace everything in development environment or when debugging is enabled
	if config.Environment == "development" || config.Debug {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	// Configure Jaeger
	if config.Instrumentation.Jaeger.Enabled {
		logger.Info("jaeger exporter enabled", nil)

		exporter, err := jaeger.NewExporter(config.Instrumentation.Jaeger.Config, errorHandler)
		emperror.Panic(err)

		trace.RegisterExporter(exporter)
	}

	// Configure graceful restart
	upg, _ := tableflip.New(tableflip.Options{})

	// Do an upgrade on SIGHUP
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGHUP)
		for range ch {
			logger.Info("graceful reloading", nil)

			_ = upg.Upgrade()
		}
	}()

	var group run.Group

	// Set up instrumentation server
	{
		const name = "instrumentation"
		logger := log.WithFields(logger, map[string]interface{}{"server": name})

		logger.Info("listening on address", map[string]interface{}{"address": config.Instrumentation.Addr})

		ln, err := upg.Fds.Listen("tcp", config.Instrumentation.Addr)
		emperror.Panic(err)

		server := &http.Server{
			Handler:  instrumentationRouter,
			ErrorLog: log.NewErrorStandardLogger(logger),
		}

		group.Add(
			func() error {
				logger.Info("starting server")

				return server.Serve(ln)
			},
			func(e error) {
				logger.Info("shutting server down")

				ctx := context.Background()
				if config.ShutdownTimeout > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, config.ShutdownTimeout)
					defer cancel()
				}

				err := server.Shutdown(ctx)
				emperror.Handle(errorHandler, emperror.With(err, "server", name))

				_ = server.Close()
			},
		)
	}

	// Register SQL stat views
	ocsql.RegisterAllViews()

	// Connect to the database
	logger.Info("connecting to database", nil)
	db, err := database.NewConnection(config.Database)
	emperror.Panic(err)
	defer db.Close()
	database.SetLogger(logger)

	// Record DB stats every 5 seconds until we exit
	defer ocsql.RecordStats(db, 5*time.Second)()

	// Register database health check
	{
		check, err := checkers.NewSQL(&checkers.SQLConfig{Pinger: db})
		emperror.Panic(errors.Wrap(err, "failed to create db health checker"))

		err = healthChecker.AddCheck(&health.Config{
			Name:     "database",
			Checker:  check,
			Interval: time.Duration(3) * time.Second,
			Fatal:    true,
		})
		emperror.Panic(errors.Wrap(err, "failed to add health checker"))
	}

	pubsub := watermill.NewPubSub(logger)
	defer pubsub.Close()

	publisher, _ := watermillx.CorrelationIDPublisherDecorator(
		watermillx.ContextCorrelationIDExtractorFunc(apptrace.CorrelationID),
	)(pubsub)

	subscriber, _ := watermillx.CorrelationIDSubscriberDecorator(
		watermillx.ContextCorrelationIDInserterFunc(apptrace.WithCorrelationID),
	)(pubsub)

	{
		h, err := watermill.NewRouter(config.Watermill.RouterConfig, logger)
		emperror.Panic(err)

		err = internal.RegisterEventHandlers(h, subscriber, logger)
		emperror.Panic(err)

		group.Add(h.Run, func(e error) { _ = h.Close() })
	}

	// Register HTTP stat views
	err = view.Register(
		ochttp.ServerRequestCountView,
		ochttp.ServerRequestBytesView,
		ochttp.ServerResponseBytesView,
		ochttp.ServerLatencyView,
		ochttp.ServerRequestCountByMethod,
		ochttp.ServerResponseCountByStatusCode,
	)
	emperror.Panic(errors.Wrap(err, "failed to register HTTP server stat views"))

	// Set up app server
	{
		const name = "app"
		logger := log.WithFields(logger, map[string]interface{}{"server": name})

		grpcServer := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))
		httpHandler, grpcHandlers := internal.NewApp(logger, publisher, errorHandler)
		httpHandler = &ochttp.Handler{
			Handler: httpHandler,
		}
		grpcHandlers(grpcServer)

		logger.Info("listening on address", map[string]interface{}{"address": config.App.Addr})

		ln, err := upg.Fds.Listen("tcp", config.App.Addr)
		emperror.Panic(err)

		server := &http.Server{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// This is a partial recreation of gRPC's internal checks:
				// https://github.com/grpc/grpc-go/blob/7346c871b018d255a1d89b3f814a645cc9c5e356/transport/handler_server.go#L61-L75
				if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
					grpcServer.ServeHTTP(w, r)
				} else {
					httpHandler.ServeHTTP(w, r)
				}
			}),
			ErrorLog: log.NewErrorStandardLogger(logger),
		}

		group.Add(
			func() error {
				logger.Info("starting server")

				return server.Serve(ln)
			},
			func(e error) {
				logger.Info("shutting server down")

				ctx := context.Background()
				if config.ShutdownTimeout > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, config.ShutdownTimeout)
					defer cancel()
				}

				err := server.Shutdown(ctx)
				emperror.Handle(errorHandler, emperror.With(err, "server", name))

				_ = server.Close()
			},
		)
	}

	// Setup signal handler
	{
		var (
			cancelInterrupt = make(chan struct{})
			ch              = make(chan os.Signal, 2)
		)
		defer close(ch)

		group.Add(
			func() error {
				signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

				select {
				case sig := <-ch:
					logger.Info("captured signal", map[string]interface{}{"signal": sig})
				case <-cancelInterrupt:
				}

				return nil
			},
			func(e error) {
				close(cancelInterrupt)
				signal.Stop(ch)
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

	err = healthChecker.Start()
	emperror.Panic(errors.Wrap(err, "failed to start health checker"))

	err = group.Run()
	emperror.Handle(errorHandler, err)
}
