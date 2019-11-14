package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/prometheus"
	"contrib.go.opencensus.io/integrations/ocsql"
	"emperror.dev/emperror"
	"emperror.dev/errors"
	logurhandler "emperror.dev/handler/logur"
	health "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	healthhttp "github.com/AppsFlyer/go-sundheit/http"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/cloudflare/tableflip"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"github.com/sagikazarmark/kitx/correlation"
	"github.com/sagikazarmark/ocmux"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
	"google.golang.org/grpc"
	invisionlog "logur.dev/integration/invision"
	"logur.dev/logur"

	"github.com/sagikazarmark/modern-go-application/internal/app/mga"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo/tododriver"
	"github.com/sagikazarmark/modern-go-application/internal/platform/buildinfo"
	"github.com/sagikazarmark/modern-go-application/internal/platform/database"
	"github.com/sagikazarmark/modern-go-application/internal/platform/log"
	"github.com/sagikazarmark/modern-go-application/internal/platform/watermill"
)

// Provisioned by ldflags
// nolint: gochecknoglobals
var (
	version    string
	commitHash string
	buildDate  string
)

func main() {
	v, p := viper.New(), pflag.NewFlagSet(friendlyAppName, pflag.ExitOnError)

	configure(v, p)

	p.String("config", "", "Configuration file")
	p.Bool("version", false, "Show version information")

	_ = p.Parse(os.Args[1:])

	if v, _ := p.GetBool("version"); v {
		fmt.Printf("%s version %s (%s) built on %s\n", friendlyAppName, version, commitHash, buildDate)

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

	err = config.Process()
	emperror.Panic(errors.WithMessage(err, "failed to process configuration"))

	// Create logger (first thing after configuration loading)
	logger := log.NewLogger(config.Log)

	// Provide some basic context to all log lines
	logger = logur.WithFields(logger, map[string]interface{}{"environment": config.Environment, "application": appName})

	log.SetStandardLogger(logger)

	if configFileNotFound {
		logger.Warn("configuration file not found")
	}

	err = config.Validate()
	if err != nil {
		logger.Error(err.Error())

		os.Exit(3)
	}

	// configure error handler
	errorHandler := logurhandler.New(logger)
	defer emperror.HandleRecover(errorHandler)

	buildInfo := buildinfo.New(version, commitHash, buildDate)

	logger.Info("starting application", buildInfo.Fields())

	telemetryRouter := http.NewServeMux()
	telemetryRouter.Handle("/version", buildinfo.Handler(buildInfo))

	// Configure health checker
	healthChecker := health.New()
	healthChecker.WithLogger(invisionlog.New(logur.WithField(logger, "component", "healthcheck")))
	{
		handler := healthhttp.HandleHealthJSON(healthChecker)
		telemetryRouter.Handle("/healthz", handler)

		// Kubernetes style health checks
		telemetryRouter.HandleFunc("/healthz/live", func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("ok"))
		})
		telemetryRouter.Handle("/healthz/ready", handler)
	}

	zpages.Handle(telemetryRouter, "/debug")

	trace.ApplyConfig(config.Opencensus.Trace.Config())

	// Configure OpenCensus exporter
	if config.Opencensus.Exporter.Enabled {
		logger.Info("opencensus exporter enabled")

		exporter, err := ocagent.NewExporter(append(
			config.Opencensus.Exporter.Options(),
			ocagent.WithServiceName(appName),
		)...)
		emperror.Panic(err)

		trace.RegisterExporter(exporter)
		view.RegisterExporter(exporter)
	}

	// Configure Prometheus exporter
	if config.Opencensus.Prometheus.Enabled {
		logger.Info("prometheus exporter enabled")

		exporter, err := prometheus.NewExporter(prometheus.Options{
			OnError: emperror.WithDetails(
				errorHandler,
				"component", "opencensus",
				"exporter", "prometheus",
			).Handle,
		})
		emperror.Panic(err)

		view.RegisterExporter(exporter)
		telemetryRouter.Handle("/metrics", exporter)
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

	// Set up telemetry server
	{
		const name = "telemetry"
		logger := logur.WithField(logger, "server", name)

		logger.Info("listening on address", map[string]interface{}{"address": config.Telemetry.Addr})

		ln, err := upg.Fds.Listen("tcp", config.Telemetry.Addr)
		emperror.Panic(err)

		server := &http.Server{
			Handler:  telemetryRouter,
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
				emperror.Handle(errorHandler, errors.WithDetails(err, "server", name))

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
	_ = healthChecker.RegisterCheck(&health.Config{
		Check:           checks.Must(checks.NewPingCheck("db.check", db, time.Millisecond*100)),
		ExecutionPeriod: 3 * time.Second,
	})

	publisher, subscriber := watermill.NewPubSub(logger)
	defer publisher.Close()
	defer subscriber.Close()

	publisher, _ = message.MessageTransformPublisherDecorator(func(msg *message.Message) {
		if cid, ok := correlation.FromContext(msg.Context()); ok {
			middleware.SetCorrelationID(cid, msg)
		}
	})(publisher)

	subscriber, _ = message.MessageTransformSubscriberDecorator(func(msg *message.Message) {
		if cid := middleware.MessageCorrelationID(msg); cid != "" {
			msg.SetContext(correlation.ToContext(msg.Context(), cid))
		}
	})(subscriber)

	{
		h, err := watermill.NewRouter(config.Watermill.RouterConfig, logger)
		emperror.Panic(err)

		err = mga.RegisterEventHandlers(h, subscriber, logger)
		emperror.Panic(err)

		group.Add(func() error { return h.Run(context.Background()) }, func(e error) { _ = h.Close() })
	}

	// Register stat views
	err = view.Register(
		// Health checks
		health.ViewCheckCountByNameAndStatus,
		health.ViewCheckStatusByName,
		health.ViewCheckExecutionTime,

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
	emperror.Panic(errors.Wrap(err, "failed to register stat views"))

	// Set up app server
	{
		const name = "app"
		logger := logur.WithField(logger, "server", name)

		httpRouter := mux.NewRouter()
		httpRouter.Use(ocmux.Middleware())

		httpServer := &http.Server{
			Handler: &ochttp.Handler{
				Handler: httpRouter,
				StartOptions: trace.StartOptions{
					Sampler:  trace.AlwaysSample(),
					SpanKind: trace.SpanKindServer,
				},
			},
			ErrorLog: log.NewErrorStandardLogger(logger),
		}

		grpcServer := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{
			StartOptions: trace.StartOptions{
				Sampler:  trace.AlwaysSample(),
				SpanKind: trace.SpanKindServer,
			},
		}))

		// In larger apps, this should be split up into smaller functions
		mga.InitializeApp(httpRouter, grpcServer, publisher, logger, errorHandler)

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
				emperror.Handle(errorHandler, errors.WithDetails(err, "server", name))

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

	err = group.Run()
	emperror.Handle(errorHandler, err)
}
