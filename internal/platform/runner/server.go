// Package netrunner provides a group runner for servers.
package runner

import (
	"context"
	"net"
	"time"

	"github.com/pkg/errors"
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

// ServerRunner implements server group run functions.
type ServerRunner struct {
	Server   Server
	Listener net.Listener

	ShutdownTimeout time.Duration

	Logger       Logger
	ErrorHandler ErrorHandler
}

// NewServerRunner returns a new ServerRunner.
func NewServerRunner(server Server, listener net.Listener) *ServerRunner {
	return &ServerRunner{
		Server:   server,
		Listener: listener,
	}
}

// Start starts the server and waits for it to return.
func (r *ServerRunner) Start() error {
	if r.Server == nil {
		return errors.New("server is not configured")
	}

	if r.Listener == nil {
		return errors.New("listener is not configured")
	}

	r.logger().Info(
		"starting server",
		map[string]interface{}{
			"network": r.Listener.Addr().Network(),
			"address": r.Listener.Addr().String(),
		},
	)

	return r.Server.Serve(r.Listener)
}

// Stop tries to shut the server down gracefully first (if the server supports it), then forcefully closes it.
func (r *ServerRunner) Stop(e error) {
	r.logger().Info("shutting server down")

	if server, ok := r.Server.(gracefulServer); ok {
		r.logger().Info("attempting graceful shutdown")

		ctx := context.Background()
		if r.ShutdownTimeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(context.Background(), r.ShutdownTimeout)

			defer cancel()
		}

		err := server.Shutdown(ctx)
		if err != nil {
			r.errorHandler().Handle(err)
		}
	}

	// TODO: add error handling
	r.Server.Close()
}

func (r *ServerRunner) logger() Logger {
	if r.Logger == nil {
		return defaultLogger
	}

	return r.Logger
}

func (r *ServerRunner) errorHandler() ErrorHandler {
	if r.ErrorHandler == nil {
		return defaultErrorHandler
	}

	return r.ErrorHandler
}
