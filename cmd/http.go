package main

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
)

// startHTTPServer starts an HTTP server.
func startHTTPServer(server *http.Server, errChan chan<- error, logger log.Logger) {
	level.Info(logger).Log(
		"msg", "starting server",
		"addr", server.Addr,
	)

	errChan <- server.ListenAndServe()
}

// stopHTTPServer starts an HTTP server.
func stopHTTPServer(ctx context.Context, server *http.Server, errorHandler emperror.Handler) {
	err := server.Shutdown(ctx)
	if err != nil {
		errorHandler.Handle(err)
	}
}
