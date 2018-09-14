package main

import (
	"context"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/conf"
	"github.com/goph/emperror"
	errorlog "github.com/goph/emperror/log"
	"github.com/pkg/errors"
	"github.com/sagikazarmark/go-service-project-boilerplate/app"
	"github.com/sagikazarmark/go-service-project-boilerplate/internal"
	"github.com/sagikazarmark/go-service-project-boilerplate/internal/helloworld"
	"github.com/sagikazarmark/go-service-project-boilerplate/internal/helloworld/driver/web"
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
	logger, err := app.NewLogger(config.Log, appCtx)
	if err != nil {
		panic(err)
	}

	// Configure error handler
	errorHandler := errorlog.NewHandler(logger)

	defer emperror.HandleRecover(errorHandler)

	// Connect to the database
	level.Debug(logger).Log("msg", "connecting to database")
	db, err := app.NewDatabaseConnection(config.Database)
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

	httpErrorLog := stdlog.New(log.NewStdlibAdapter(level.Error(logger)), "", 0)

	httpServer := &http.Server{
		Addr:     config.HTTPAddr,
		Handler:  router,
		ErrorLog: httpErrorLog,
	}
	defer httpServer.Close()

	level.Info(logger).Log(
		"version", appCtx.Build.Version,
		"commit_hash", appCtx.Build.CommitHash,
		"build_date", appCtx.Build.Date,
		"msg", "starting",
	)

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
