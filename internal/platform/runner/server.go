package runner

import (
	"context"
	"net"
	"time"
)

type server interface {
	Serve(l net.Listener) error
	Shutdown(ctx context.Context) error
	Close() error
}

// logger is the fundamental interface for all log operations.
type logger interface {
	// Info logs an info event.
	Info(msg ...interface{})
}

// errorHandler is responsible for handling an error.
type errorHandler interface {
	// Handle takes care of unhandled errors.
	Handle(err error)
}

// Server implements server group run functions.
type Server struct {
	Server   server
	Listener net.Listener

	ShutdownTimeout time.Duration

	Logger       logger
	ErrorHandler errorHandler
}

// Start starts the server and waits for it to return.
func (r *Server) Start() error {
	r.Logger.Info("starting server")

	return r.Server.Serve(r.Listener)
}

// Stop tries to shut the server down gracefully first, then forcefully closes it.
func (r *Server) Stop(e error) {
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
