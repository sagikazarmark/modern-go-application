package runner

import (
	"context"
	"net"
	"time"
)

// Server is the underlying server instance.
type Server interface {
	Serve(l net.Listener) error
	Shutdown(ctx context.Context) error
	Close() error
}

// ServerRunnerLogger is the fundamental interface for all log operations.
type ServerRunnerLogger interface {
	// Info logs an info event.
	Info(msg string, fields ...map[string]interface{})
}

// ServerRunnerErrorHandler is responsible for handling an error.
type ServerRunnerErrorHandler interface {
	// Handle takes care of unhandled errors.
	Handle(err error)
}

// ServerRunner implements server group run functions.
type ServerRunner struct {
	Server   Server
	Listener net.Listener

	ShutdownTimeout time.Duration

	Logger       ServerRunnerLogger
	ErrorHandler ServerRunnerErrorHandler
}

// Start starts the server and waits for it to return.
func (r *ServerRunner) Start() error {
	r.Logger.Info("starting server", nil)

	return r.Server.Serve(r.Listener)
}

// Stop tries to shut the server down gracefully first, then forcefully closes it.
func (r *ServerRunner) Stop(e error) {
	ctx := context.Background()
	if r.ShutdownTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), r.ShutdownTimeout)

		defer cancel()
	}

	r.Logger.Info("shutting server down")

	err := r.Server.Shutdown(ctx)
	if err != nil {
		r.ErrorHandler.Handle(err)
	}

	r.Server.Close()
}
