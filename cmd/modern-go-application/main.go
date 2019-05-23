package main

import (
	"context"
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
	"github.com/sagikazarmark/modern-go-application/internal/platform/watermill"
	"github.com/sagikazarmark/modern-go-application/internal/todo/tododriver"
	"github.com/sagikazarmark/modern-go-application/pkg/correlation"
)

// Provisioned by ldflags
// nolint: gochecknoglobals
var (
	version    string
	commitHash string
	buildDate  string
)

func main() {
	v, p := viper.New(), pflag.NewFlagSet(firendlyAppName, pflag.ExitOnError)

	configure(v, p)

	p.String("config", "", "Configuration file")
	p.Bool("version", false, "Show version information")
	p.Bool("dump-config", false, "Dump configuration to the console (and exit)")

	_ = p.Parse(os.Args[1:])

	if v, _ := p.GetBool("version"); v {
		fmt.Printf("%s version %s (%s) built on %s\n", firendlyAppName, version, commitHash, buildDate)

		os.Exit(0)
	}

	if c, _ := p.GetString("config"); c != "" {
		v.SetConfigFile(c)
	}

	err := v.ReadInConfig()
	_, configFileNotFound := err.(viper.ConfigFileNotFoundError)
	if !configFileNotFound {
		emperror.Panic(errors.Wrap(err, "failed to read configuration"))
	}

	var config configuration
	err = v.Unmarshal(&config)
	emperror.Panic(errors.Wrap(err, "failed to unmarshal configuration"))

	// Create logger (first thing after configuration loading)
	logger := log.NewLogger(config.Log)

	// Provide some basic context to all log lines
	logger = log.WithFields(logger, map[string]interface{}{"environment": config.Environment, "application": appName})

	log.SetStandardLogger(logger)

	if configFileNotFound {
		logger.Warn("configuration file not found")
	}

	err = config.Validate()
	if err != nil {
		logger.Error(err.Error())

		os.Exit(3)
	}

	if d, _ := p.GetBool("dump-config"); d {
		fmt.Printf("%+v\n", config)

		os.Exit(0)
	}

	// configure error handler
	errorHandler := errorhandler.New(logger)
	defer emperror.HandleRecover(errorHandler)

	buildInfo := buildinfo.New(version, commitHash, buildDate)

	logger.Info("starting application", buildInfo.Fields())

	instrumentationRouter := http.NewServeMux()
	instrumentationRouter.Handle("/version", buildinfo.Handler(buildInfo))

	// configure health checker
	healthChecker := healthcheck.New(logger)
	instrumentationRouter.Handle("/healthz", healthcheck.Handler(healthChecker))

	trace.ApplyConfig(config.Opencensus.Trace.Config())

	// configure Prometheus
	if config.Instrumentation.Prometheus.Enabled {
		logger.Info("prometheus exporter enabled")

		exporter, err := prometheus.NewExporter(config.Instrumentation.Prometheus.Config, errorHandler)
		emperror.Panic(err)

		view.RegisterExporter(exporter)
		instrumentationRouter.Handle("/metrics", exporter)
	}

	// configure Jaeger
	if config.Instrumentation.Jaeger.Enabled {
		logger.Info("jaeger exporter enabled")

		exporter, err := jaeger.NewExporter(config.Instrumentation.Jaeger.Config, errorHandler)
		emperror.Panic(err)

		trace.RegisterExporter(exporter)
	}

	// configure graceful restart
	upg, _ := tableflip.New(tableflip.Options{})

	// Do an upgrade on SIGHUP
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGHUP)
		for range ch {
			logger.Info("graceful reloading")

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
	logger.Info("connecting to database")
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
		watermillx.ContextCorrelationIDExtractorFunc(correlation.ID),
	)(pubsub)

	subscriber, _ := watermillx.CorrelationIDSubscriberDecorator(
		watermillx.ContextCorrelationIDInserterFunc(correlation.WithID),
	)(pubsub)

	{
		h, err := watermill.NewRouter(config.Watermill.RouterConfig, logger)
		emperror.Panic(err)

		err = internal.RegisterEventHandlers(h, subscriber, logger)
		emperror.Panic(err)

		group.Add(h.Run, func(e error) { _ = h.Close() })
	}

	// Register stat views
	err = view.Register(
		// HTTP
		ochttp.ServerRequestCountView,
		ochttp.ServerRequestBytesView,
		ochttp.ServerResponseBytesView,
		ochttp.ServerLatencyView,
		ochttp.ServerRequestCountByMethod,
		ochttp.ServerResponseCountByStatusCode,

		// GRPC
		ocgrpc.ServerReceivedBytesPerRPCView,
		ocgrpc.ServerSentBytesPerRPCView,
		ocgrpc.ServerLatencyView,
		ocgrpc.ServerCompletedRPCsView,

		// Todo
		tododriver.CreatedTodoCountView,
		tododriver.DoneTodoCountView,
	)
	emperror.Panic(errors.Wrap(err, "failed to register HTTP server stat views"))

	// Set up app server
	{
		const name = "app"
		logger := log.WithFields(logger, map[string]interface{}{"server": name})

		grpcServer := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{
			StartOptions: trace.StartOptions{
				Sampler:  trace.AlwaysSample(),
				SpanKind: trace.SpanKindServer,
			},
		}))

		httpHandler, grpcHandlers := internal.NewApp(logger, publisher, errorHandler)
		httpHandler = &ochttp.Handler{
			Handler: httpHandler,
			StartOptions: trace.StartOptions{
				Sampler:  trace.AlwaysSample(),
				SpanKind: trace.SpanKindServer,
			},
		}
		grpcHandlers(grpcServer)

		httpServer := &http.Server{
			Handler:  httpHandler,
			ErrorLog: log.NewErrorStandardLogger(logger),
		}

		logger.Info("listening on address", map[string]interface{}{"address": config.App.HttpAddr})

		httpLn, err := upg.Fds.Listen("tcp", config.App.HttpAddr)
		emperror.Panic(err)

		logger.Info("listening on address", map[string]interface{}{"address": config.App.GrpcAddr})

		grpcLn, err := upg.Fds.Listen("tcp", config.App.GrpcAddr)
		emperror.Panic(err)

		group.Add(
			func() error {
				logger.Info("starting server")

				return httpServer.Serve(httpLn)
			},
			func(e error) {
				logger.Info("shutting server down")

				ctx := context.Background()
				if config.ShutdownTimeout > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, config.ShutdownTimeout)
					defer cancel()
				}

				err := httpServer.Shutdown(ctx)
				emperror.Handle(errorHandler, emperror.With(err, "server", name))

				_ = httpServer.Close()
			},
		)

		group.Add(
			func() error {
				logger.Info("starting server")

				return grpcServer.Serve(grpcLn)
			},
			func(e error) {
				logger.Info("shutting server down")

				defer grpcServer.Stop()
				grpcServer.GracefulStop()
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
