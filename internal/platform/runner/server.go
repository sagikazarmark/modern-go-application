package runner

import (
	"context"
	"net"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
)

type server interface {
	Serve(l net.Listener) error
	Shutdown(ctx context.Context) error
	Close() error
}

// Server implements server group run functions.
type Server struct {
	Server   server
	Listener net.Listener

	ShutdownTimeout time.Duration

	Logger       log.Logger
	ErrorHandler emperror.Handler
}

// Start starts the server and waits for it to return.
func (r *Server) Start() error {
	level.Info(r.Logger).Log("msg", "starting server")

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

	level.Info(r.Logger).Log("msg", "shutting server down")

	err := r.Server.Shutdown(ctx)
	if err != nil {
		r.ErrorHandler.Handle(err)
	}

	r.Server.Close()
}
