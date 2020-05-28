package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	_ "net/http/pprof" // register pprof HTTP handlers #nosec
	"os"
	"os/signal"
	"syscall"
	"time"

	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/prometheus"
	"contrib.go.opencensus.io/integrations/ocsql"
	"emperror.dev/emperror"
	"emperror.dev/errors"
	"emperror.dev/errors/match"
	logurhandler "emperror.dev/handler/logur"
	health "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	healthhttp "github.com/AppsFlyer/go-sundheit/http"
	"github.com/cloudflare/tableflip"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"github.com/sagikazarmark/appkit/buildinfo"
	appkiterrors "github.com/sagikazarmark/appkit/errors"
	appkitrun "github.com/sagikazarmark/appkit/run"
	"github.com/sagikazarmark/ocmux"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
	"google.golang.org/grpc"
	"logur.dev/logur"

	"github.com/sagikazarmark/modern-go-application/internal/app/mga"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo/tododriver"
	"github.com/sagikazarmark/modern-go-application/internal/common/commonadapter"
	"github.com/sagikazarmark/modern-go-application/internal/platform/appkit"
	"github.com/sagikazarmark/modern-go-application/internal/platform/database"
	"github.com/sagikazarmark/modern-go-application/internal/platform/gosundheit"
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

const (
	// appName is an identifier-like name used anywhere this app needs to be identified.
	//
	// It identifies the application itself, the actual instance needs to be identified via environment
	// and other details.
	appName = "mga"

	// friendlyAppName is the visible name of the application.
	friendlyAppName = "Modern Go Application"
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

	// Override the global standard library logger to make sure everything uses our logger
	log.SetStandardLogger(logger)

	if configFileNotFound {
		logger.Warn("configuration file not found")
	}

	err = config.Validate()
	if err != nil {
		logger.Error(err.Error())

		os.Exit(3)
	}

	// Configure error handler
	errorHandler := logurhandler.New(logger)
	defer emperror.HandleRecover(errorHandler)

	buildInfo := buildinfo.New(version, commitHash, buildDate)

	logger.Info("starting application", buildInfo.Fields())

	telemetryRouter := http.NewServeMux()
	telemetryRouter.Handle("/buildinfo", buildinfo.HTTPHandler(buildInfo))

	// Register pprof endpoints
	telemetryRouter.Handle("/debug/pprof/", http.DefaultServeMux)

	// Configure health checker
	healthChecker := health.New()
	healthChecker.WithCheckListener(gosundheit.NewLogger(logur.WithField(logger, "component", "healthcheck")))
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
		exporter, err := ocagent.NewExporter(append(
			config.Opencensus.Exporter.Options(),
			ocagent.WithServiceName(appName),
		)...)
		emperror.Panic(err)

		trace.RegisterExporter(exporter)
		view.RegisterExporter(exporter)
	}

	// Configure Prometheus exporter
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
		defer server.Close()

		group.Add(
			func() error { return server.Serve(ln) },
			func(err error) { _ = server.Shutdown(context.Background()) },
		)
	}

	// Register SQL stat views
	ocsql.RegisterAllViews()

	// Connect to the database
	logger.Info("connecting to database")
	dbConnector, err := database.NewConnector(config.Database)
	emperror.Panic(err)

	database.SetLogger(logger)

	db := sql.OpenDB(dbConnector)
	defer db.Close()

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

	publisher = watermill.PublisherCorrelationID(publisher)
	subscriber = watermill.SubscriberCorrelationID(subscriber)

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
		tododriver.CreatedTodoItemCountView,
		tododriver.CompleteTodoItemCountView,
	)
	emperror.Panic(errors.Wrap(err, "failed to register stat views"))

	// Set up app server
	{
		const name = "app"
		logger := logur.WithField(logger, "server", name)

		httpRouter := mux.NewRouter()
		httpRouter.Use(ocmux.Middleware())

		cors := handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete}),
			handlers.AllowedHeaders([]string{"content-type"}),
		)

		httpServer := &http.Server{
			Handler: &ochttp.Handler{
				// Handler: httpRouter,
				Handler: cors(httpRouter),
				StartOptions: trace.StartOptions{
					Sampler:  trace.AlwaysSample(),
					SpanKind: trace.SpanKindServer,
				},
				IsPublicEndpoint: true,
			},
			ErrorLog: log.NewErrorStandardLogger(logger),
		}
		defer httpServer.Close()

		grpcServer := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{
			StartOptions: trace.StartOptions{
				Sampler:  trace.AlwaysSample(),
				SpanKind: trace.SpanKindServer,
			},
			IsPublicEndpoint: true,
		}))
		defer grpcServer.Stop()

		// In larger apps, this should be split up into smaller functions
		{
			logger := commonadapter.NewContextAwareLogger(logger, appkit.ContextExtractor)
			errorHandler := emperror.WithFilter(
				emperror.WithContextExtractor(errorHandler, appkit.ContextExtractor),
				appkiterrors.IsServiceError, // filter out service errors
			)

			mga.InitializeApp(httpRouter, grpcServer, publisher, config.App.Storage, db, logger, errorHandler)

			h, err := watermill.NewRouter(logger)
			emperror.Panic(err)

			err = mga.RegisterEventHandlers(h, subscriber, logger)
			emperror.Panic(err)

			group.Add(func() error { return h.Run(context.Background()) }, func(e error) { _ = h.Close() })
		}

		logger.Info("listening on address", map[string]interface{}{"address": config.App.HttpAddr})

		httpLn, err := upg.Fds.Listen("tcp", config.App.HttpAddr)
		emperror.Panic(err)

		logger.Info("listening on address", map[string]interface{}{"address": config.App.GrpcAddr})

		grpcLn, err := upg.Fds.Listen("tcp", config.App.GrpcAddr)
		emperror.Panic(err)

		group.Add(
			func() error { return httpServer.Serve(httpLn) },
			func(err error) { _ = httpServer.Shutdown(context.Background()) },
		)
		group.Add(
			func() error { return grpcServer.Serve(grpcLn) },
			func(err error) { grpcServer.GracefulStop() },
		)
	}

	// Setup signal handler
	group.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	// Setup graceful restart
	group.Add(appkitrun.GracefulRestart(context.Background(), upg))

	err = group.Run()
	emperror.WithFilter(errorHandler, match.As(&run.SignalError{}).MatchError).Handle(err)
}
