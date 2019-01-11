// Package netrunner provides a group runner for servers.
package netrunner

import (
	"context"
	"net"
	"time"
)

// Server accepts connections from the network and answers to requests.
type Server interface {
	Serve(l net.Listener) error
	Close() error
}

// gracefulServer allows to try graceful shutdown first.
type gracefulServer interface {
	Shutdown(ctx context.Context) error
}

// Logger is the fundamental interface for all log operations.
type Logger interface {
	// Info logs an info event.
	Info(msg string, fields ...map[string]interface{})
}

// ErrorHandler is responsible for handling an error.
type ErrorHandler interface {
	// Handle takes care of unhandled errors.
	Handle(err error)
}

// ServerRunner implements server group run functions.
type ServerRunner struct {
	Server   Server
	Listener net.Listener

	ShutdownTimeout time.Duration

	Logger       Logger
	ErrorHandler ErrorHandler
}

// Start starts the server and waits for it to return.
func (r *ServerRunner) Start() error {
	r.Logger.Info(
		"starting server",
		map[string]interface{}{
			"network": r.Listener.Addr().Network(),
			"address":    r.Listener.Addr().String(),
		},
	)

	return r.Server.Serve(r.Listener)
}

// Stop tries to shut the server down gracefully first (if the server supports it), then forcefully closes it.
func (r *ServerRunner) Stop(e error) {
	r.Logger.Info("shutting server down")

	if server, ok := r.Server.(gracefulServer); ok {
		r.Logger.Info("attempting graceful shutdown")

		ctx := context.Background()
		if r.ShutdownTimeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(context.Background(), r.ShutdownTimeout)

			defer cancel()
		}

		err := server.Shutdown(ctx)
		if err != nil {
			r.ErrorHandler.Handle(err)
		}
	}

	// TODO: add error handling (when the returned error is an unexpected one and not eg. http.ErrServerClosed)
	r.Server.Close()
}
