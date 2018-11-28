package runner

import (
	"context"
	"net"
	"time"

	"github.com/goph/emperror"
)

type server interface {
	Serve(l net.Listener) error
	Shutdown(ctx context.Context) error
	Close() error
}

// logger is the fundamental interface for all log operations.
type logger interface {
	// Infof logs an info event and optionally formats the message.
	Infof(msg string, args ...interface{})
}

// Server implements server group run functions.
type Server struct {
	Server   server
	Listener net.Listener

	ShutdownTimeout time.Duration

	Logger       logger
	ErrorHandler emperror.Handler
}

// Start starts the server and waits for it to return.
func (r *Server) Start() error {
	r.Logger.Infof("starting server")

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

	r.Logger.Infof("shutting server down")

	err := r.Server.Shutdown(ctx)
	if err != nil {
		r.ErrorHandler.Handle(err)
	}

	r.Server.Close()
}
